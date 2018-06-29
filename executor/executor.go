package executor

import (
	"log"
)

type Executor struct {
	maxWorkers    int
	ActiveWorkers int

	Jobs    chan Job
	Reports chan Report
	signals chan int
}

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
				e.AddReport(r)
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

func (e *Executor) launchWorker(job Job, reports chan<- Report) {
	report := job.Execute()

	if len(reports) < cap(reports) {
		reports <- report
	} else {
		log.Println("Executor's report channel is full...")
	}

	e.ActiveWorkers--
}

func (e *Executor) AddJob(job Job) bool {
	if len(e.Jobs) == cap(e.Jobs) {

		return false
	}

	e.Jobs <- job
	return true
}

func (e *Executor) AddReport(report Report) bool {
	if len(e.Reports) == cap(e.Reports) {

		return false
	}

	e.Reports <- report
	return true
}

func (e *Executor) inactive() bool {
	return e.ActiveWorkers == 0 && len(e.Jobs) == 0 && len(e.Reports) == 0
}
