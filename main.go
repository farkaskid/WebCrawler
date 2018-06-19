package main

import (
	"WebCrawler/crawler"
	"flag"
	"log"
	"sync"
)

func main() {
	url := flag.String("url", "", "url to start")
	boundedBy := flag.String("bound", "", "Domain to bound the crawler")

	flag.Parse()

	channel := make(chan bool, 1)

	crawler := crawler.Crawler{
		Processor: crawler.LogProcessor{},
		Collector: crawler.URLCollector{make(map[string]bool), make(map[string]bool), &sync.Mutex{}},
		Filter:    crawler.CrossDomainFilter{*boundedBy},
		Url:       *url,
		Done:      channel,
	}

	crawler.Done <- false
	go crawler.Start()

	for activeCrawlers := 1; activeCrawlers >= 1; {
		select {
		case status := <-channel:
			if !status {
				activeCrawlers++
			} else {
				activeCrawlers--

				if activeCrawlers == 1 {
					activeCrawlers--
				}
			}
		}
	}

	log.Println("Crawler finished")
}
