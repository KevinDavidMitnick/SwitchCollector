package service

import (
	"github.com/SwitchCollector/service/device"
	"github.com/SwitchCollector/service/flow"
	"github.com/SwitchCollector/service/visit"
)

func Init() {
	device.Init()
}

func CollectFlow() {
	go visit.Init()
	go visit.CleanStale()
	go flow.Collect()
	go flow.CleanStale()
	select {}
}
