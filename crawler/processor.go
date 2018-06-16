package crawler

import (
	"log"
)

type LogProcessor struct {
	Url string
}

func (processor LogProcessor) Process(data []string) {
	log.Println(len(data), "URLs found on the URL", processor.Url)
}
