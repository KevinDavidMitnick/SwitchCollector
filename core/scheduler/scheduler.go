package scheduler

import (
	"fmt"
	"sync"
	"time"
)

//Object interface ,scheduler run
type Object interface {
	Run(timestamp int64)
}

//Scheduler struct ,scheduler queue
type Scheduler struct {
	Queue map[int64][]Object `json:"queue"`
	sync.RWMutex
}

//Run scheduler
func (scheduler *Scheduler) Run(interval int64) {
	timer := time.NewTicker(time.Second * time.Duration(interval))
	for {
		select {
		case <-timer.C:
			scheduler.RLock()
			timestamp := time.Now().Unix()
			if len(scheduler.Queue[interval]) == 0 {
				scheduler.RUnlock()
				break
			}
			for _, obj := range scheduler.Queue[interval] {
				go obj.Run(timestamp)
			}
			scheduler.RUnlock()
		case <-time.After(time.Second * time.Duration(interval*2)):
			fmt.Println("timeout for scheduler...")
		}
	}
}

//Scheduler ,to run task
func (scheduler *Scheduler) Scheduler() {
	scheduler.RLock()
	defer scheduler.RUnlock()
	for interval := range scheduler.Queue {
		go scheduler.Run(interval)
	}
}

//GetScheduler , get scheduler
func GetScheduler() *Scheduler {
	var scheduler Scheduler
	scheduler.Queue = make(map[int64][]Object)
	return &scheduler
}
