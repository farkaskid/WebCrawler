package crawler

import (
	"net/url"
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

	for _, rawurl := range urls {
		url, err := url.Parse(rawurl)

		if err != nil {
			continue
		}

		if url.Hostname() == filter.Domain {
			filteredURLs = append(filteredURLs, rawurl)
		}
	}

	return filteredURLs
}
