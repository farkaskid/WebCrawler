package crawler

import (
	"WebCrawler/executor"
	"fmt"
	"net/http"
)

type CrawlReport struct {
	URL           string
	HttpStatus    int
	Err           string
	ConnectedURLs []string
}

func (report CrawlReport) Status() int {
	return report.HttpStatus
}

func (report CrawlReport) String() string {
	if report.HttpStatus == 0 {
		return "Failed to crawl URL: " + report.URL + ". Cause: " + report.Err
	}

	return fmt.Sprintf("Found %d URLs on %s which responded with %d. ", len(report.ConnectedURLs), report.URL, report.HttpStatus)
}

type DefaultProcessor struct{}

func (processor DefaultProcessor) Process(requestedURL string, res *http.Response, connectedURLs []string, err error) executor.Report {
	if res == nil {
		return CrawlReport{requestedURL, 0, err.Error(), make([]string, 0)}
	}

	return CrawlReport{requestedURL, res.StatusCode, "", connectedURLs}
}
