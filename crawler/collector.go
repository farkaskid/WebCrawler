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
		log.Println(rawurl, ": is visited.")

		return rawurls
	}

	res, err := http.Get(rawurl)

	if err != nil {
		log.Println("Failed to crawl URL", rawurl)

		return rawurls
	}

	if 200 > res.StatusCode || res.StatusCode >= 400 {
		log.Println("Bad response", res.StatusCode, "on:", rawurl)

		return rawurls
	}

	redirectedUrl := res.Request.URL.String()

	if collector.Visited(existingUrls, redirectedUrl) {
		log.Println(redirectedUrl, ": is visited.")

		return rawurls
	}

	collector.Add(existingUrls, redirectedUrl, true)
	collector.Add(existingUrls, rawurl, true)

	defer res.Body.Close()

	content := readResponse(res.Body)

	urlGen := urlGenerator(string(content))

	for childRawUrl := urlGen(); childRawUrl != ""; childRawUrl = urlGen() {
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

func urlGenerator(allContent string) func() string {
	content := allContent

	return func() string {
		start := strings.Index(content, "href=\"")

		if start == -1 {
			return ""
		}

		content = content[start+6:]
		end := strings.Index(content, "\"")

		if end == -1 {
			return ""
		}

		url := content[:end]
		content = content[end:]

		return url
	}
}
