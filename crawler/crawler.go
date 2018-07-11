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

func (crawler *Crawler) spawnChild(rawurl string) {
	child := Crawler{
		URL:       rawurl,
		Processor: crawler.Processor,
		Filter:    crawler.Filter,
		Collector: crawler.Collector,
		Executor:  crawler.Executor,
	}

	crawler.Executor.AddTask(Task{child})
}

type Task struct {
	Crawler Crawler
}

func (task Task) Execute() executor.Report {
	crawler := task.Crawler

	response, URLs, err := crawler.Collect(crawler.URL)
	report, count := crawler.Process(crawler.URL, response, URLs, err), 0

	for _, url := range crawler.Filter.Filter(URLs) {
		crawler.spawnChild(url)
		count++
	}

	return report
}

func (task Task) String() string {
	return "Crawl url: " + task.Crawler.URL
}
