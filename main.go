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
	rawurl := flag.String("url", "", "url to start.")
	bound := flag.Bool("bound", false, "Domain to bound the crawler.")
	maxWorkers := flag.Int("maxWorkers", 1000, "Number of concurrent crawler tasks.")

	flag.Parse()

	var filter crawler.Filter

	url, err := url.Parse(*rawurl)

	if err != nil {
		log.Fatalln(err)
	}

	if *bound {
		filter = &crawler.CrossDomainFilter{url.Hostname()}
	} else {
		filter = &crawler.NoneFilter{}
	}

	ctlCh := make(chan int)

	exec := executor.NewExecutor(*maxWorkers, ctlCh)
	reports := exec.Reports
	jobs := exec.Jobs

	c := newCrawler(*rawurl, filter, &exec)

	exec.AddJob(crawler.CrawlerJob{c})

	defer log.Println("Crawler finished")

	for {
		select {
		case <-reports:
			if len(jobs) == 0 && len(reports) == 0 && exec.ActiveWorkers == 0 {
				log.Println("Sending termination request...")
				ctlCh <- 1
				if 0 == <-ctlCh {
					return
				}
			}
		}
	}
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
