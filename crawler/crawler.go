package crawler

import (
	"WebCrawler/executor"
)

type Crawler struct {
	Url string

	Collector
	Processor
	Filter
	*executor.Executor
}

func (crawler *Crawler) spawnChild(resource string) {
	child := Crawler{
		Url:       resource,
		Processor: crawler.Processor,
		Filter:    crawler.Filter,
		Collector: crawler.Collector,
		Executor:  crawler.Executor,
	}

	crawler.Executor.Add(CrawlerJob{child})
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

	return CrawlerReport{count}
}

type CrawlerReport struct {
	urlCount int
}

func (report CrawlerReport) Status() int {
	return 0
}
