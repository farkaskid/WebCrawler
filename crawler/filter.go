package crawler

import (
	"strings"
)

type NoneFilter struct{}

func (filter NoneFilter) Filter(urls []string) []string {
	return urls
}

type CrossDomainFilter struct {
	Domain string
}

func (filter CrossDomainFilter) Filter(urls []string) []string {
	var filteredURLs []string

	for _, url := range urls {
		if strings.Contains(url, filter.Domain) {
			filteredURLs = append(filteredURLs, url)
		}
	}

	return filteredURLs
}
