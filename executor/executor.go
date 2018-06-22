package executor

import (
	"sync"
)

type Executor struct {
	maxWorkers    int
	ActiveWorkers int

	Jobs    chan Job
	Reports chan Report
	ctl     <-chan int

	mutex *sync.Mutex
}

func NewExecutor(maxWorkers int, ctl <-chan int) Executor {
	e := Executor{
		maxWorkers: maxWorkers,
		Jobs:       make(chan Job, maxWorkers),
		Reports:    make(chan Report, maxWorkers),
		ctl:        ctl,
		mutex:      &sync.Mutex{},
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
				if signal == 1 {
					return 0
				}

			case j := <-e.Jobs:
				if e.ActiveWorkers < e.maxWorkers {
					e.ActiveWorkers++
					go func() {
						reports <- j.Execute()
						e.ActiveWorkers--
					}()
				}

			case r := <-reports:
				if r.Status() == 0 {
					e.Reports <- r
				}
			}
		}
	}()
}

func (e Executor) Add(job Job) bool {
	if len(e.Jobs) == e.maxWorkers {
		return false
	}

	e.Jobs <- job

	return true
}
