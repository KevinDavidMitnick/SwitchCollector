package g

import "sync"

type DataValue struct {
	LastValue interface{} `json:"lastvalue"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

type MetricData struct {
	MetricType string                  `json:"merictype"`
	Data       map[string][]*DataValue `json:"data"`
}

type DeviceData struct {
	Metrics map[string]map[string]*MetricData `json:"metrics"`
}

var (
	locker       sync.RWMutex
	globalData   DeviceData
	indexNameMap map[string]string
)

func GetGlobalData() *DeviceData {
	locker.Lock()
	defer locker.Unlock()

	if globalData.Metrics == nil {
		globalData.Metrics = make(map[string]map[string]*MetricData)
	}
	return &globalData
}

func GetIndexNameMap() map[string]string {
	locker.Lock()
	defer locker.Unlock()

	if indexNameMap == nil {
		indexNameMap = make(map[string]string)
	}

	return indexNameMap
}
