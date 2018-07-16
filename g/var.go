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

type DeviceList struct {
	Data []string `json:"Data"`
}

var (
	locker       sync.RWMutex
	globalData   DeviceData
	indexNameMap map[string]map[string]string
)

func GetGlobalData() *DeviceData {
	locker.Lock()
	defer locker.Unlock()

	if globalData.Metrics == nil {
		globalData.Metrics = make(map[string]map[string]*MetricData)
	}
	return &globalData
}

func GetIndexNameMap() map[string]map[string]string {
	locker.Lock()
	defer locker.Unlock()

	if indexNameMap == nil {
		indexNameMap = make(map[string]map[string]string)
	}

	return indexNameMap
}

func GetDeviceList() *DeviceList {
	var ret DeviceList
	ret.Data = make([]string, 0)
	for ip := range indexNameMap {
		ret.Data = append(ret.Data, ip)
	}
	return &ret
}
