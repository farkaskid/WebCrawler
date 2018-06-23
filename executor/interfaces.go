package executor

type Report interface {
	Status() int
	String() string
}

type Job interface {
	Execute() Report
	String() string
}
