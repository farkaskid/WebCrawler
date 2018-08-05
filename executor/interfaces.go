package executor

// Report interface defines a report of a Job. A report must define how the Job went using the
// Status() method and a the description about the report using the String().
type Report interface {
	Status() int
	String() string
}

// Task interface defines a task for the Executor. 'Execute' method is used to start the task which
// must return a Report. String() gives a small description about the Task.
type Task interface {
	Execute() Report
	String() string
}
