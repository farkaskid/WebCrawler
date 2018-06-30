// Package executor contains a job executor and its related utilities.
// A Job executor provides a simple interface for concurrently conducting multiple small number
// of tasks using a single goroutine for each one. It offers a signals channel to control the executor
// itself and a reports channel to get reports of the jobs as they get completed.
// The Job and Report are both interfaces and the client can specify any kind of logic in them.
// But it is recommended that a Job should contain a task that is simple and individual and should NOT
// spawn further goroutines.
// User can control the maximum simultaneous goroutines.
package executor
