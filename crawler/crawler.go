package crawler

import (
	"WebCrawler/executor"
)

type Crawler struct {
	Url string

	Collector
	Processor
	Filter
	executor.Executor
}

type CrawlerReport struct {
	urlCount int
}

func (report CrawlerReport) Status() int {
	return 0
}

type CrawlerJob struct {
	Crawler Crawler
}

func (job CrawlerJob) Execute() executor.Report {
	crawler := job.Crawler
	data := crawler.Collect(crawler.Url)

	// crawler.Process(crawler.data)

	count := 0
	for _, datum := range crawler.Filter.Filter(data) {
		crawler.spawnChild(datum)
		count++
	}

	return CrawlerReport{count}
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
