package crawler

import (
	"bytes"
	"errors"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
)

type URLCollector struct {
	UrlMap map[uint64]bool
	*sync.Mutex
}

func (collector *URLCollector) sync(f func()) {
	collector.Lock()
	f()
	collector.Unlock()
}

func (collector *URLCollector) Visited(m map[uint64]bool, s string) (visited bool) {
	collector.sync(func() { visited = m[hash(s)] })
	return
}

func (collector *URLCollector) Present(m map[uint64]bool, s string) (present bool) {
	collector.sync(func() { _, present = m[hash(s)] })
	return
}

func (collector *URLCollector) Add(m map[uint64]bool, s string, visited bool) {
	collector.sync(func() { m[hash(s)] = visited })
}

func convertToAbs(parentUrl *url.URL, childUrl *url.URL) string {
	parentUrl.Path = path.Join(parentUrl.Path, childUrl.Path)

	return parentUrl.String()
}

func (collector *URLCollector) Collect(rawurl string) []string {
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

	content := readContent(res.Body)

	urlGen := urlGenerator(content)

	for childRawUrl, err := urlGen(); err == nil; childRawUrl, err = urlGen() {
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

func hash(s string) uint64 {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))

	if err != nil {
		log.Println("Failed to hash value:", s)
		return 0
	}

	return h.Sum64()
}

func readContent(readerCloser io.ReadCloser) string {
	var buf bytes.Buffer

	io.Copy(&buf, readerCloser)
	readerCloser.Close()

	return buf.String()
}

func urlGenerator(allContent string) func() (string, error) {
	content := allContent

	return func() (string, error) {
		start := strings.Index(content, "href=\"")

		if start == -1 {
			return "", errors.New("content exhausted.")
		}

		content = content[start+6:]
		end := strings.Index(content, "\"")

		if end == -1 {
			return "", nil
		}

		url := content[:end]
		content = content[end:]

		return url, nil
	}
}
