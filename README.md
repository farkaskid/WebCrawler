# WebCrawler

A simple web crawler. This leverages Go's lightweight goroutines to achieve high level of concurrency and thus showing impressive performance. Just try

`crawl -url http://www.google.com`

![](demo.gif)

and enjoy.

OSX binary is available in the release. You can download it and run with the usual `./crawlerOSX` syntax.

## More Options

- `-bound` flag can be used to bound the crawler within the domain of the given URL.
`crawl -url http://www.netflix.com -bound`
- `-maxWorkers` flag can be used to control the maximum number of simultaneous goroutines the crawler will spawn.
`crawl -url http://www.netflix.com -maxWorkers 10000`
It defaults to 1000.

## Reporting

The crawler uses simple files to report the crawling results. A `CrawlReport` is generated for each crawled URL which contains:
- URL: The URL whose report is this.
- HttpStatus: If HTTP GET request was successful then the status of the response.
- Err: Error if any with the technical details of failure.
- ConnectedURLs: A list of URLs that were found in the response body of this URL.

All the reports can be found in the `reports` folder. A single report file contains a bunch of CrawlReports encoded in the specified encoding. Currently, `gob` and `json` encoding are supported, `gob` being the default. You can pass the `-json` flag to encode the CrawlReports in JSON encoding.

The number of CrawlReports that a report file will contain can be changed by passing the `-reportSize` flag(ex `-reportSize 300`), it defaults to 500. It should be noted that setting a very high value can cause the crawler to consume large amounts of memory.

Example: `crawl -url http://www.netflix.com -bound -json -reportSize 100`
