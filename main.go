package main

import (
	"WebCrawler/crawler"
	"flag"
	"sync"
)

func main() {
	url := flag.String("url", "", "url to start")

	flag.Parse()

	channel := make(chan bool)

	crawler := crawler.Crawler{
		Processor: crawler.LogProcessor{},
		Collector: crawler.URLCollector{make(map[string]bool), make(map[string]bool), &sync.Mutex{}},
		Url:       *url,
		Done:      channel,
	}

	go crawler.Start()

	for {
		select {
		case <-channel:
		}
	}
}
