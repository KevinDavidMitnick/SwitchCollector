package scheduler

import (
	"fmt"
	"time"
)

//Object interface ,scheduler run
type Object interface {
	Run()
}

//Scheduler struct ,scheduler queue
type Scheduler struct {
	Queue map[int64][]Object `json:"queue"`
}

//runScheduler ,run scheduler after interval
func (scheduler *Scheduler) run(interval int64, tasks []Object) {
	for {
		timer := time.NewTicker(time.Second * time.Duration(interval))
		select {
		case <-timer.C:
			for _, obj := range tasks {
				go obj.Run()
			}
		case <-time.After(time.Second * time.Duration(interval*2)):
			fmt.Println("timeout for scheduler...")
		}
	}
}

//Scheduler ,to run task
func (scheduler *Scheduler) Scheduler() {
	for interval, schedulers := range scheduler.Queue {
		go scheduler.run(interval, schedulers)
	}
}

//GetScheduler, get scheduler
func GetScheduler() *Scheduler {
	var scheduler Scheduler
	scheduler.Queue = make(map[int64][]Object)
	return &scheduler
}
