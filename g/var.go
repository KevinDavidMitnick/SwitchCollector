package g

import (
	"sync"
	"time"
)

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
	Data           []string `json:"Data"`
	StatisticsTime int64    `json:"StatisticsTime"`
}

type DeviceInfo struct {
	Data           []map[string]interface{} `json:"Data"`
	StatisticsTime int64                    `json:"StatisticsTime"`
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
	ret.StatisticsTime = time.Now().Unix()
	return &ret
}

func GetDeviceInfo(ip string) *DeviceInfo {
	var ret DeviceInfo
	ret.Data = make([]map[string]interface{}, 0)
	for name, metricData := range globalData.Metrics["ip"] {
		info := make(map[string]interface{})
		if metricData.MetricType == "infos" || metricData.MetricType == "metrics" {
			info[name] = metricData.Data["liucong"][0].Value
			ret.Data = append(ret.Data, info)
		}
	}
	ret.StatisticsTime = time.Now().Unix()
	return &ret
}
