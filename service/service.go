package service

import (
	"github.com/SwitchCollector/service/device"
)

func CollectFlow() {
	device := device.GetDevice()
	device.InitTasks()
	device.InitScheduler()
	device.Collect()
	go device.CleanStale()
	go device.UpdateScheduler()
}
