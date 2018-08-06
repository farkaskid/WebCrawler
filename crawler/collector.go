package crawler

import (
	"bytes"
	"errors"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// URLCollector is an implementation of the collector interface that collects Anchors for the pages
// that the crawler visits.
type URLCollector struct {
	URLMap      map[uint64]bool
	AnchorRegex *regexp.Regexp
	*sync.Mutex
}

func (collector *URLCollector) sync(f func()) {
	collector.Lock()
	f()
	collector.Unlock()
}

func (collector *URLCollector) visited(s string) (visited bool) {
	collector.sync(func() { visited = collector.URLMap[hash(s)] })
	return
}

func (collector *URLCollector) present(s string) (present bool) {
	collector.sync(func() { _, present = collector.URLMap[hash(s)] })
	return
}

func (collector *URLCollector) add(s string, visited bool) {
	collector.sync(func() { collector.URLMap[hash(s)] = visited })
}

// Collect method collects all the URLs that are available in the form on href="" markup.
func (collector *URLCollector) Collect(URL *url.URL) (*http.Response, []Anchor, error) {
	var anchors []Anchor
	rawurl := URL.String()

	if collector.visited(rawurl) {
		return nil, anchors, errors.New("URL is already crawled")
	}

	res, err := http.Get(rawurl)

	if err != nil {
		return res, anchors, err
	}

	if 200 > res.StatusCode || res.StatusCode >= 400 {
		return res, anchors, errors.New("URL responded with status code " + res.Status)
	}

	pageurl := res.Request.URL.String()

	if collector.visited(pageurl) {
		return res, anchors, errors.New("URL is already crawled")
	}

	collector.add(pageurl, true)
	collector.add(rawurl, true)

	for _, anchor := range collector.findAnchors(readContent(res.Body)) {
		if len(anchor.Href) < 5 {
			continue
		}

		URL, err := url.Parse(anchor.Href)

		if err != nil {
			continue
		}

		if !URL.IsAbs() {
			pageURL, _ := url.Parse(pageurl)
			anchor.Href = pageURL.ResolveReference(URL).String()
		}

		if collector.present(anchor.Href) {
			continue
		}

		collector.add(anchor.Href, false)
		anchors = append(anchors, anchor)
	}

	return res, anchors, nil
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

func (collector *URLCollector) findAnchors(text string) []Anchor {
	var anchors []Anchor

	for _, match := range collector.AnchorRegex.FindAllString(text, -1) {
		anchor, err := newAnchor(match)

		if err == nil {
			anchors = append(anchors, anchor)
		}
	}

	return anchors
}

func newAnchor(text string) (Anchor, error) {
	var anchor Anchor
	var start, end int

	start = strings.Index(text, `href="`)

	if start == -1 {
		return anchor, errors.New("Failed to extract href")
	}

	start += 6
	end = start + strings.Index(text[start:], `"`)
	anchor.Href = text[start:end]

	start = strings.Index(text, ">")

	if start == -1 {
		return anchor, errors.New("Failed to extract title")
	}

	start++
	end = start + strings.Index(text[start:], "<")
	anchor.Title = text[start:end]

	return anchor, nil
}
