package main

import (
	"WebCrawler/crawler"
	"WebCrawler/executor"
	"flag"
	"log"
	"net/url"
	"sync"
)

func main() {
	rawurl := flag.String("url", "", "url to start")
	bound := flag.Bool("bound", false, "Domain to bound the crawler")

	flag.Parse()

	var filter crawler.Filter

	url, err := url.Parse(*rawurl)

	if err != nil {
		log.Fatalln(err)
	}

	if *bound {
		log.Println(url.Hostname())
		filter = &crawler.CrossDomainFilter{url.Hostname()}
	} else {
		filter = &crawler.NoneFilter{}
	}

	ctlCh := make(chan int)

	exec := executor.NewExecutor(100000, ctlCh)
	reports := exec.Reports
	jobs := exec.Jobs

	c := newCrawler(*rawurl, filter, &exec)

	exec.Add(crawler.CrawlerJob{c})

	for {
		select {
		case <-reports:
			if len(jobs) == 0 && len(reports) != 0 && exec.ActiveWorkers == 0 {
				break
			}
		}
	}

	log.Println("Crawler finished")
}

func newCrawler(url string, filter crawler.Filter, executor *executor.Executor) crawler.Crawler {
	processor := &crawler.LogProcessor{}
	collector := &crawler.URLCollector{make(map[string]bool), &sync.Mutex{}}

	return crawler.Crawler{
		Processor: processor,
		Collector: collector,
		Filter:    filter,
		Url:       url,
		Executor:  executor,
	}
}
