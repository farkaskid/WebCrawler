package executor

import (
	"log"
)

type Executor struct {
	maxWorkers    int
	ActiveWorkers int

	Jobs    chan Job
	Reports chan Report
	ctl     chan int
}

func NewExecutor(maxWorkers int, ctl chan int) Executor {
	chanSize := 1000

	if maxWorkers > chanSize {
		chanSize = maxWorkers
	}

	e := Executor{
		maxWorkers: maxWorkers,
		Jobs:       make(chan Job, chanSize),
		Reports:    make(chan Report, chanSize),
		ctl:        ctl,
	}

	e.init()

	return e
}

func (e Executor) init() {
	reports := make(chan Report, e.maxWorkers)

	go func() int {
		for {
			select {
			case signal := <-e.ctl:
				log.Println("Received termination request...")
				if signal == 1 {
					if e.inactive() {
						log.Println("No active jobs, exiting...")
						e.ctl <- 0
						return 0
					}

					e.ctl <- 1
					log.Println("Some jobs are still active...")
				}

			case j := <-e.Jobs:
				if e.ActiveWorkers < e.maxWorkers {
					e.ActiveWorkers++

					go func() {
						report := j.Execute()

						if len(reports) < cap(reports) {
							reports <- report
						} else {
							log.Println("Executor's report channel is full...")
						}

						e.ActiveWorkers--
					}()
				} else {
					e.AddJob(j)
				}

			case r := <-reports:
				if r.Status() == 0 {
					e.AddReport(r)
				}
			}
		}
	}()
}

func (e Executor) AddJob(job Job) bool {
	if len(e.Jobs) == cap(e.Jobs) {

		return false
	}

	e.Jobs <- job
	return true
}

func (e Executor) AddReport(report Report) bool {
	if len(e.Reports) == cap(e.Reports) {

		return false
	}

	e.Reports <- report
	return true
}

func (e Executor) inactive() bool {
	return e.ActiveWorkers == 0 && len(e.Jobs) == 0 && len(e.Reports) == 0
}
