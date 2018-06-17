package crawler

type Collector interface {
	Collect(resource string) []string
}

type Processor interface {
	Process(data []string)
}

type Crawler struct {
	Url  string
	Done chan<- bool

	data []string

	Collector
	Processor
}

func (crawler *Crawler) Start() {
	crawler.data = crawler.Collect(crawler.Url)

	// crawler.Process(crawler.data)

	for _, datum := range crawler.data {
		crawler.spawnChild(datum)
	}

	crawler.Done <- true
}

func (crawler *Crawler) spawnChild(resource string) {
	child := Crawler{
		Url:       resource,
		Done:      crawler.Done,
		Processor: crawler.Processor,

		Collector: crawler.Collector,
	}

	child.Done <- false
	go child.Start()
}
