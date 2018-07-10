package crawler

import (
	"WebCrawler/executor"
	"net/http"
)

type Collector interface {
	Collect(resource string) (*http.Response, []string, error)
}

type Processor interface {
	Process(requestedURL string, response *http.Response, connectedURLs []string, err error) executor.Report
}

type Filter interface {
	Filter(urls []string) []string
}
