package crawler

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type URLCollector struct {
	ResourceMap map[string]bool
	UrlMap      map[string]bool
	Mutex       *sync.Mutex
}

func (collector URLCollector) Collect(url string) []string {
	var urls []string

	collector.Mutex.Lock()
	if collector.UrlMap[url] {
		collector.Mutex.Unlock()
		return urls
	}
	collector.Mutex.Unlock()

	res, err := http.Get(url)

	if err != nil {
		log.Println("Failed to crawl URL", url)

		return urls
	}

	defer res.Body.Close()

	content := readResponse(res.Body)

	collector.Mutex.Lock()
	collector.UrlMap[url] = true
	collector.Mutex.Unlock()

	urlFinder := urlFinderGenerator(string(content))

	for childURL := urlFinder(); childURL != ""; childURL = urlFinder() {
		if len(childURL) < 5 {
			continue
		}

		collector.Mutex.Lock()
		if _, contains := collector.ResourceMap[childURL]; contains {
			collector.Mutex.Unlock()
			continue
		}
		collector.Mutex.Unlock()

		if childURL[:4] != "http" {
			collector.Mutex.Lock()
			collector.ResourceMap[childURL] = false
			collector.Mutex.Unlock()
			childURL = url + childURL
		}

		collector.Mutex.Lock()
		if _, contains := collector.UrlMap[childURL]; contains {
			collector.Mutex.Unlock()
			continue
		}
		collector.Mutex.Unlock()

		collector.Mutex.Lock()
		collector.UrlMap[childURL] = false
		collector.Mutex.Unlock()

		urls = append(urls, childURL)
	}

	log.Println(len(urls), "URLs found on the URL", url)

	return urls
}

func readResponse(reader io.Reader) []byte {
	var content []byte
	buffer := make([]byte, 1024)

	for c, err := reader.Read(buffer); c == 1024 && err == nil; c, err = reader.Read(buffer) {
		content = append(content, buffer...)
	}

	return content
}

func urlFinderGenerator(content string) func() string {
	modifiedContent := content

	return func() string {
		start := strings.Index(modifiedContent, "href=\"")

		if start == -1 {
			return ""
		}

		modifiedContent = modifiedContent[start+6:]
		end := strings.Index(modifiedContent, "\"")

		if end == -1 {
			return ""
		}

		url := modifiedContent[:end]

		modifiedContent = modifiedContent[end:]

		return url
	}
}
