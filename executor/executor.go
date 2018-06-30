package executor

import (
	"log"
)

// The Executor struct is the main executor for jobs.
// 'maxWorkers' represents the maximum number of simultaneous goroutines.
// 'ActiveWorkers' tells the number of active goroutines spawned by the Executor at given time.
// 'Jobs' is the channel on which the Executor receives the jobs.
// 'Reports' is channel on which the Executor publishes the every jobs reports.
// 'signals' is channel that can be used to control the executor. Right now, only the termination
// signal is supported which is essentially is sending '1' on this channel by the client.
type Executor struct {
	maxWorkers    int
	ActiveWorkers int

	Jobs    chan Job
	Reports chan Report
	signals chan int
}

// NewExecutor creates a new Executor.
// 'maxWorkers' tells the maximum number of simultaneous goroutines.
// 'signals' channel can be used to control the Executor.
func NewExecutor(maxWorkers int, signals chan int) *Executor {
	chanSize := 1000

	if maxWorkers > chanSize {
		chanSize = maxWorkers
	}

	e := Executor{
		maxWorkers: maxWorkers,
		Jobs:       make(chan Job, chanSize),
		Reports:    make(chan Report, chanSize),
		signals:    signals,
	}

	go e.launch()

	return &e
}

// launch starts the main loop for polling on the all the relevant channels and handling differents
// messages.
func (e *Executor) launch() int {
	reports := make(chan Report, e.maxWorkers)

	for {
		select {
		case signal := <-e.signals:
			if e.handleSignals(signal) == 0 {
				return 0
			}

		case r := <-reports:
			if r.Status() == 0 {
				e.addReport(r)
			}

		default:
			if e.ActiveWorkers < e.maxWorkers && len(e.Jobs) > 0 {
				j := <-e.Jobs
				e.ActiveWorkers++
				go e.launchWorker(j, reports)
			}
		}
	}
}

// handleSignals is called whenever anything is received on the 'signals' channel.
// It performs the relevant task according to the received signal(request) and then responds either
// with 0 or 1 indicating whether the request was respected(0) or rejected(1).
func (e *Executor) handleSignals(signal int) int {
	if signal == 1 {
		log.Println("Received termination request...")

		if e.inactive() {
			log.Println("No active workers, exiting...")
			e.signals <- 0
			return 0
		}

		e.signals <- 1
		log.Println("Some jobs are still active...")
	}

	return 1
}

// launchWorker is called whenever a new Job is received and Executor can spawn more workers to spawn
// a new Worker.
// Each worker is launched on a new goroutine. It processes the given job and publishes the report on
// the Executor's internal reports channel.
func (e *Executor) launchWorker(job Job, reports chan<- Report) {
	report := job.Execute()

	if len(reports) < cap(reports) {
		reports <- report
	} else {
		log.Println("Executor's report channel is full...")
	}

	e.ActiveWorkers--
}

// AddJob is used to submit a new job to the Executor is a non-blocking way. The Client can submit
// a new job using the Executor's Jobs channel directly but that will block if the Jobs channel is
// full.
// It should be considered that this method doesn't add the given job if the Jobs channel is full
// and it is up to client to try again later.
func (e *Executor) AddJob(job Job) bool {
	if len(e.Jobs) == cap(e.Jobs) {

		return false
	}

	e.Jobs <- job
	return true
}

// addReport is used by the Executor to publish the reports in a non-blocking way. It client is not
// reading the reports channel or is slower that the Executor publishing the reports, the Executor's
// reports channel is going to get full. In that case this method will not block and that report will
// not be added.
func (e *Executor) addReport(report Report) bool {
	if len(e.Reports) == cap(e.Reports) {

		return false
	}

	e.Reports <- report
	return true
}

// inactive checks if the Executor is idle. This happens when there are no pending jobs, active
// workers and reports to publish. It is called when a the Executor receives a termination request.
func (e *Executor) inactive() bool {
	return e.ActiveWorkers == 0 && len(e.Jobs) == 0 && len(e.Reports) == 0
}
