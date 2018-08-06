package crawler

import (
	"WebCrawler/executor"
	"net/http"
	"net/url"
)

// Anchor represent an HTML anchor tag.
type Anchor struct {
	Href  string
	Title string
}

type Collector interface {
	Collect(url *url.URL) (*http.Response, []Anchor, error)
}

type Processor interface {
	Process(url *url.URL, response *http.Response, connectedURLs []Anchor, err error) executor.Report
}

type Filter interface {
	Filter(urls []string) []string
}
