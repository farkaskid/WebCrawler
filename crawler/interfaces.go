package crawler

type Collector interface {
	Collect(resource string) []string
}

type Processor interface {
	Process(data []string)
}

type Filter interface {
	Filter(urls []string) []string
}
