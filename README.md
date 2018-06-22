# WebCrawler

A simple web crawler. This leverages Go's lightweight goroutines to achieve high level of concurrency and thus showing impressive performance. Just try

`crawl -url http://www.google.com`

and enjoy.

OSX binary is available in the release. You can download it and run with the usual `./crawlerOSX` syntax.

## More Options

- `-bound` flag can be used to bound the crawler within the domain of the given URL.
`crawl -url http://www.netflix.com -bound`
