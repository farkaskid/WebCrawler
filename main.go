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
		filter = crawler.CrossDomainFilter{url.Hostname()}
	} else {
		filter = crawler.NoneFilter{}
	}

	ctlCh := make(chan int)

	exec := executor.NewExecutor(10, ctlCh)
	reports := exec.Reports
	jobs := exec.Jobs

	c := crawler.Crawler{
		Processor: crawler.LogProcessor{},
		Collector: crawler.URLCollector{make(map[string]bool), make(map[string]bool), &sync.Mutex{}},
		Filter:    filter,
		Url:       *rawurl,
		Executor:  exec,
	}

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
