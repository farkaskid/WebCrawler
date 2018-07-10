package crawler

import (
	"WebCrawler/executor"
	"fmt"
	"net/http"
)

// CrawlReport represents the results of the crawling task of a single URL.
type CrawlReport struct {
	URL           string
	HTTPStatus    int
	Err           string
	ConnectedURLs []string
}

// Status returns the HTTPStatus of the GET request on the URL of this report.
func (report CrawlReport) Status() int {
	return report.HTTPStatus
}

func (report CrawlReport) String() string {
	if report.HTTPStatus == 0 {
		return "Failed to crawl URL: " + report.URL + ". Cause: " + report.Err
	}

	return fmt.Sprintf("Found %d URLs on %s which responded with %d. ", len(report.ConnectedURLs),
		report.URL, report.HTTPStatus)
}

// DefaultProcessor is the default implementation of the crawler.Processor interface. This just
// creates a appropriate CrawlReport instance from the results of the collector.
type DefaultProcessor struct{}

// Process creates a CrawlReport instance from the given parameters. It should be noted that when
// http.Response is nil then the HTTPStatus in the CrawlReport is set as 0.
func (processor DefaultProcessor) Process(requestedURL string, res *http.Response, connectedURLs []string,
	err error) executor.Report {
	if res == nil {
		return CrawlReport{requestedURL, 0, err.Error(), make([]string, 0)}
	}

	return CrawlReport{requestedURL, res.StatusCode, "", connectedURLs}
}
