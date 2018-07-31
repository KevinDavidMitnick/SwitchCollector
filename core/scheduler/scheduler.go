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

//Scheduler ,to run task
func (scheduler *Scheduler) Scheduler() {
	scheduler.RLock()
	defer scheduler.RUnlock()
	for interval := range scheduler.Queue {
		go func() {
			for {
				timer := time.NewTicker(time.Second * time.Duration(interval))
				select {
				case <-timer.C:
					timestamp := time.Now().Unix()
					if len(scheduler.Queue[interval]) == 0 {
						break
					}
					for _, obj := range scheduler.Queue[interval] {
						go obj.Run(timestamp)
					}
				case <-time.After(time.Second * time.Duration(interval*2)):
					fmt.Println("timeout for scheduler...")
				}
			}
		}()
	}
}

//GetScheduler, get scheduler
func GetScheduler() *Scheduler {
	var scheduler Scheduler
	scheduler.Queue = make(map[int64][]Object)
	return &scheduler
}
