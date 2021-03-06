package main

import (
	"WebCrawler/crawler"
	"WebCrawler/executor"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"sync"
)

var (
	rawurl string

	bound  bool
	isJSON bool

	maxWorkers       int
	reportSize       int
	reportFileNumber int

	reportsBuf []executor.Report
)

func init() {
	flag.StringVar(&rawurl, "url", "", "url to start.")

	flag.BoolVar(&bound, "bound", false, "Domain to bound the crawler.")
	flag.BoolVar(&isJSON, "json", false, "Should generate reports in JSON..?")

	flag.IntVar(&maxWorkers, "maxWorkers", 1000, "Number of concurrent crawler tasks.")
	flag.IntVar(&reportSize, "reportSize", 500, "Size of a single report.")

	flag.Parse()

	os.Mkdir("reports", os.ModePerm)

	gob.Register(crawler.Report{})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func newCrawler(executor *executor.Executor) crawler.Crawler {
	URL, err := url.Parse(rawurl)
	check(err)
	URL.Fragment = ""

	var filter crawler.Filter

	if bound {
		filter = &crawler.CrossDomainFilter{URL.Hostname()}
	} else {
		filter = &crawler.NoneFilter{}
	}

	processor := &crawler.DefaultProcessor{}

	regex, _ := regexp.Compile("<a[^>]*>([^<]+)</a>")
	collector := &crawler.URLCollector{make(map[uint64]bool), regex, &sync.Mutex{}}

	return crawler.Crawler{
		Processor: processor,
		Collector: collector,
		Filter:    filter,
		URL:       URL,
		Executor:  executor,
	}
}

func handleReport(report executor.Report) {
	log.Println(report)
	if report.Status() == 0 {
		return
	}

	reportsBuf = append(reportsBuf, report)

	if len(reportsBuf) == reportSize {
		writeReports()
		reportFileNumber++
		reportsBuf = reportsBuf[:0]
	}
}

func writeReports() {
	var buf bytes.Buffer
	var err error
	var extension string

	if isJSON {
		encoder := json.NewEncoder(&buf)
		err, extension = encoder.Encode(reportsBuf), "json"
	} else {
		encoder := gob.NewEncoder(&buf)
		err, extension = encoder.Encode(reportsBuf), "gob"
	}

	check(err)
	err = ioutil.WriteFile(fmt.Sprintf("reports/#%d.%s", reportFileNumber, extension), buf.Bytes(), 0644)
	check(err)
}

func shutdownExecutor(signals chan int) bool {
	log.Println("Sending termination request...")
	signals <- 1

	if 0 == <-signals {
		writeReports()
		return true
	}

	return false
}

func main() {
	signals := make(chan int)
	e := executor.NewExecutor(maxWorkers, signals)
	reports := e.Reports

	c := newCrawler(e)

	e.AddTask(crawler.Task{c})

	for {
		select {
		case r := <-reports:
			handleReport(r)

			if e.Inactive() {
				if shutdownExecutor(signals) {
					log.Println("Crawler finished")

					return
				}
			}
		}
	}
}
