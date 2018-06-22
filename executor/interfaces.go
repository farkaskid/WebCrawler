package executor

type Report interface {
	Status() int
}

type Job interface {
	Execute() Report
}
