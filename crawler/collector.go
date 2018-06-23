package crawler

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
)

type URLCollector struct {
	UrlMap map[string]bool
	*sync.Mutex
}

func (collector URLCollector) sync(f func()) {
	collector.Lock()
	f()
	collector.Unlock()
}

func (collector URLCollector) Visited(m map[string]bool, s string) (visited bool) {
	collector.sync(func() { visited = m[s] })
	return
}

func (collector URLCollector) Present(m map[string]bool, s string) (present bool) {
	collector.sync(func() { _, present = m[s] })
	return
}

func (collector URLCollector) Add(m map[string]bool, s string, visited bool) {
	collector.sync(func() { m[s] = visited })
}

func convertToAbs(parentUrl *url.URL, childUrl *url.URL) string {
	parentUrl.Path = path.Join(parentUrl.Path, childUrl.Path)

	return parentUrl.String()
}

func (collector URLCollector) Collect(rawurl string) []string {
	var rawurls []string

	_, err := url.Parse(rawurl)

	if err != nil {
		return rawurls
	}

	existingUrls := collector.UrlMap

	if collector.Visited(existingUrls, rawurl) {
		return rawurls
	}

	res, err := http.Get(rawurl)

	if err != nil {
		log.Println("Failed to crawl URL", rawurl)

		return rawurls
	}

	redirectedUrl := res.Request.URL.String()

	if err == nil {
		if collector.Visited(existingUrls, redirectedUrl) {
			return rawurls
		}
	}

	collector.Add(existingUrls, redirectedUrl, true)
	collector.Add(existingUrls, rawurl, true)

	defer res.Body.Close()

	content := readResponse(res.Body)

	urlFinder := urlFinderGenerator(string(content))

	for childRawUrl := urlFinder(); childRawUrl != ""; childRawUrl = urlFinder() {
		if len(childRawUrl) < 5 {
			continue
		}

		childurl, err := url.Parse(childRawUrl)

		if err != nil {
			continue
		}

		if !childurl.IsAbs() {
			RedirectedUrl, _ := url.Parse(redirectedUrl)

			childRawUrl = convertToAbs(RedirectedUrl, childurl)
		}

		if collector.Present(existingUrls, childRawUrl) {
			continue
		}

		collector.Add(existingUrls, childRawUrl, false)

		rawurls = append(rawurls, childRawUrl)
	}

	if len(rawurls) >= 0 {
		log.Println(len(rawurls), "URLs found on the URL", redirectedUrl)
	}

	return rawurls
}

func readResponse(reader io.Reader) []byte {
	var content []byte
	buffer := make([]byte, 1024)

	for c, err := reader.Read(buffer); c > 0 || err == nil; c, err = reader.Read(buffer) {
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
