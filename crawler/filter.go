package crawler

import (
	"net/url"
)

type NoneFilter struct{}

func (filter NoneFilter) Filter(urls []string) []string {
	var filteredURLs []string

	for _, u := range urls {
		if u == "" {
			continue
		}

		if _, err := url.Parse(u); err != nil {
			continue
		}

		filteredURLs = append(filteredURLs, u)
	}

	return filteredURLs
}

type CrossDomainFilter struct {
	Domain string
}

func (filter CrossDomainFilter) Filter(urls []string) []string {
	var filteredURLs []string

	for _, rawurl := range urls {
		if rawurl == "" {
			continue
		}

		URL, err := url.Parse(rawurl)

		if err != nil {
			continue
		}

		if URL.Hostname() == filter.Domain {
			filteredURLs = append(filteredURLs, rawurl)
		}
	}

	return filteredURLs
}
