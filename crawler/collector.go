package crawler

import (
	"bytes"
	"errors"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// URLCollector is an implementation of the collector interface that collects URLs for the pages
// that the crawler visits.
type URLCollector struct {
	URLMap map[uint64]bool
	*sync.Mutex
}

func (collector *URLCollector) sync(f func()) {
	collector.Lock()
	f()
	collector.Unlock()
}

func (collector *URLCollector) visited(m map[uint64]bool, s string) (visited bool) {
	collector.sync(func() { visited = m[hash(s)] })
	return
}

func (collector *URLCollector) present(m map[uint64]bool, s string) (present bool) {
	collector.sync(func() { _, present = m[hash(s)] })
	return
}

func (collector *URLCollector) add(m map[uint64]bool, s string, visited bool) {
	collector.sync(func() { m[hash(s)] = visited })
}

// Collect method collects all the URLs that are available in the form on href="" markup.
func (collector *URLCollector) Collect(rawurl string) []string {
	var rawurls []string

	_, err := url.Parse(rawurl)

	if err != nil {
		return rawurls
	}

	existingURLs := collector.URLMap

	if collector.visited(existingURLs, rawurl) {
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

	redirectedURL := res.Request.URL.String()

	if collector.visited(existingURLs, redirectedURL) {
		log.Println(redirectedURL, ": is visited.")

		return rawurls
	}

	collector.add(existingURLs, redirectedURL, true)
	collector.add(existingURLs, rawurl, true)

	content := readContent(res.Body)

	urlGen := urlGenerator(content)

	for childRawURL, err := urlGen(); err == nil; childRawURL, err = urlGen() {
		if len(childRawURL) < 5 {
			continue
		}

		childURL, err := url.Parse(childRawURL)

		if err != nil {
			continue
		}

		if !childURL.IsAbs() {
			RedirectedURL, _ := url.Parse(redirectedURL)

			childRawURL = RedirectedURL.ResolveReference(childURL).String()
		}

		if collector.present(existingURLs, childRawURL) {
			continue
		}

		collector.add(existingURLs, childRawURL, false)

		rawurls = append(rawurls, childRawURL)
	}

	if len(rawurls) >= 0 {
		log.Println(len(rawurls), "URLs found on the URL", redirectedURL)
	}

	return rawurls
}

// This function is a utility used to hash the given string using the FNV-1a hashing algorithm.
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
			return "", errors.New("content exhausted")
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
