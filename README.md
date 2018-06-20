# WebCrawler

A simple web crawler. This leverages Go's lightweight goroutines to achieve high level of concurrency and thus showing impressive performance. Just try

`crawl -url http://www.google.com`

and enjoy.

OSX binary is available in the release.

## More Options

- `-bound` flag can be used to bound the crawler within the domain of the given URL.
`crawl -url http://www.netflix.com -bound`
