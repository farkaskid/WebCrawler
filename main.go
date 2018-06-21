package main

import (
	"WebCrawler/crawler"
	"flag"
	"log"
	"net/url"
	"sync"
)

func main() {
	rawurl := flag.String("url", "", "url to start")
	bound := flag.Bool("bound", false, "Domain to bound the crawler")
	host := flag.String("host", "", "Domain to bound the crawler")

	flag.Parse()

	var filter crawler.Filter

	url, err := url.Parse(*rawurl)
	hostName := url.Hostname()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(*bound)

	if *bound {
		if *host != "" {
			hostName = *host
		}

		log.Println(hostName)

		filter = crawler.CrossDomainFilter{hostName}
	} else {
		filter = crawler.NoneFilter{}
	}

	channel := make(chan bool, 1)

	crawler := crawler.Crawler{
		Processor: crawler.LogProcessor{},
		Collector: crawler.URLCollector{make(map[string]bool), make(map[string]bool), &sync.Mutex{}},
		Filter:    filter,
		Url:       *rawurl,
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
