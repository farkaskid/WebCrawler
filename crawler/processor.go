package crawler

import (
	"WebCrawler/executor"
	"fmt"
	"net/http"
	"net/url"
)

// Report represents the results of the crawling task of a single URL.
type Report struct {
	Url        string
	HTTPStatus int
	Err        string
	Anchors    []Anchor
}

// Status returns the HTTPStatus of the GET request on the URL of this report.
func (report Report) Status() int {
	return report.HTTPStatus
}

func (report Report) String() string {
	if report.HTTPStatus == 0 {
		return "Failed to crawl URL: " + report.Url + ". Cause: " + report.Err
	}

	return fmt.Sprintf("Found %d URLs on %s which responded with %d. ", len(report.Anchors),
		report.Url, report.HTTPStatus)
}

// DefaultProcessor is the default implementation of the crawler.Processor interface. This just
// creates a appropriate CrawlReport instance from the results of the collector.
type DefaultProcessor struct{}

// Process creates a CrawlReport instance from the given parameters. It should be noted that when
// http.Response is nil then the HTTPStatus in the CrawlReport is set as 0.
func (processor DefaultProcessor) Process(URL *url.URL, res *http.Response, anchors []Anchor,
	err error) executor.Report {
	if res == nil {
		return Report{URL.String(), 0, err.Error(), make([]Anchor, 0)}
	}

	return Report{URL.String(), res.StatusCode, "", anchors}
}
