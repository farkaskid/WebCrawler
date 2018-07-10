package crawler

import (
	"WebCrawler/executor"
)

type Crawler struct {
	URL string

	Collector
	Processor
	Filter
	Executor *executor.Executor
}

func (crawler *Crawler) spawnChild(resource string) {
	child := Crawler{
		URL:       resource,
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

	res, URLs, err := c.Collect(c.URL)
	report, count := c.Process(c.URL, res, URLs, err), 0

	for _, url := range c.Filter.Filter(URLs) {
		c.spawnChild(url)
		count++
	}

	return report
}

func (job CrawlerJob) String() string {
	return "Crawl url: " + job.Crawler.URL
}
