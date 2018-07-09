package service

import (
	"github.com/SwitchCollector/service/flow"
	"github.com/SwitchCollector/service/visit"
)

func CollectFlow() {
	go visit.Init()
	go visit.CleanStale()
	go flow.Collect()
	go flow.CleanStale()
	select {}
}
