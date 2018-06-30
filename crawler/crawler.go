package crawler

import (
	"WebCrawler/executor"
)

type Crawler struct {
	Url string

	Collector
	Processor
	Filter
	Executor *executor.Executor
}

func (crawler *Crawler) spawnChild(resource string) {
	child := Crawler{
		Url:       resource,
		Processor: crawler.Processor,
		Filter:    crawler.Filter,
		Collector: crawler.Collector,
		Executor:  crawler.Executor,
	}

	crawler.Executor.AddJob(CrawlerJob{child})
}

type CrawlerJob struct {
	Crawler Crawler
}

func (job CrawlerJob) Execute() executor.Report {
	c := job.Crawler

	urls := c.Collect(c.Url)

	// crawler.Process(crawler.data)

	count := 0
	for _, url := range c.Filter.Filter(urls) {
		c.spawnChild(url)
		count++
	}

	return CrawlerReport{c.Url, count}
}

func (job CrawlerJob) String() string {
	return "Crawl url: " + job.Crawler.Url
}

type CrawlerReport struct {
	Url      string
	UrlCount int
}

func (report CrawlerReport) Status() int {
	return 0
}

func (report CrawlerReport) String() string {
	return "Crawl report for url: " + report.Url
}
