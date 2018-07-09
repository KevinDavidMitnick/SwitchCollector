package service

import (
	"github.com/SwitchCollector/service/flow"
)

func CollectFlow() {
	go flow.Collect()
	go flow.CleanStale()
	select {}
}
