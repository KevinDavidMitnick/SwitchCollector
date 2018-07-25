package g

import (
	"strings"
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

type InterfaceInfo struct {
	Data           map[string]map[string]interface{} `json:"Data"`
	StatisticsTime int64                             `json:"StatisticsTime"`
}

type InterfaceMetric struct {
	Data           map[string]map[string][]interface{} `json:"Data"`
	StatisticsTime int64                               `json:"StatisticsTime"`
}

var (
	Locker       sync.RWMutex
	globalData   DeviceData
	indexNameMap map[string]map[string]string
)

func GetGlobalData() *DeviceData {
	Locker.Lock()
	defer Locker.Unlock()

	if globalData.Metrics == nil {
		globalData.Metrics = make(map[string]map[string]*MetricData)
	}
	return &globalData
}

func GetIndexNameMap() map[string]map[string]string {
	Locker.Lock()
	defer Locker.Unlock()

	if indexNameMap == nil {
		indexNameMap = make(map[string]map[string]string)
	}

	return indexNameMap
}

func GetDeviceList() *DeviceList {
	Locker.RLock()
	defer Locker.RUnlock()

	var ret DeviceList
	ret.Data = make([]string, 0)
	for ip := range indexNameMap {
		ret.Data = append(ret.Data, ip)
	}
	ret.StatisticsTime = time.Now().Unix()
	return &ret
}

func GetDeviceInfo(ip string) *DeviceInfo {
	Locker.RLock()
	defer Locker.RUnlock()

	var ret DeviceInfo
	ret.Data = make([]map[string]interface{}, 0)
	for metricName, metricData := range globalData.Metrics[ip] {
		info := make(map[string]interface{})
		if metricData.MetricType == "infos" || metricData.MetricType == "metrics" {
			info[metricName] = metricData.Data["liucong"][0].Value
			ret.Data = append(ret.Data, info)
		}
	}
	ret.StatisticsTime = time.Now().Unix()
	return &ret
}

func GetInterfaceInfo(ip string, filter string, accurate string) *InterfaceInfo {
	Locker.RLock()
	defer Locker.RUnlock()

	var ret InterfaceInfo
	ret.Data = make(map[string]map[string]interface{})
	for metricName, metricData := range globalData.Metrics[ip] {
		if metricData.MetricType == "multiinfos" || metricData.MetricType == "multimetrics" {
			for interfaceName, values := range metricData.Data {
				if interfaceName == "NULL0" {
					continue
				}
				if (accurate == "true" && interfaceName == filter) || (accurate != "true" && strings.Contains(interfaceName, filter)) {
					if ret.Data[interfaceName] == nil {
						ret.Data[interfaceName] = make(map[string]interface{})
					}
					length := len(values)
					ret.Data[interfaceName][metricName] = values[length-1].Value
				}
			}
		}
	}

	ret.StatisticsTime = time.Now().Unix()
	return &ret
}

func GetInterfaceMetric(ip string, filter string, accurate string, period int64) *InterfaceMetric {
	Locker.RLock()
	defer Locker.RUnlock()

	var ret InterfaceMetric
	ret.Data = make(map[string]map[string][]interface{})
	startTime := time.Now().Unix() - period
	if period == 0 {
		startTime = 0
	}
	for metricName, metricData := range globalData.Metrics[ip] {
		if metricData.MetricType == "multimetrics" {
			for interfaceName, values := range metricData.Data {
				if interfaceName == "NULL0" {
					continue
				}
				if (accurate == "true" && interfaceName == filter) || (accurate != "true" && strings.Contains(interfaceName, filter)) {
					if ret.Data[interfaceName] == nil {
						ret.Data[interfaceName] = make(map[string][]interface{})
					}
					length := len(values)
					for i := length - 1; i >= 0; i-- {
						if values[i].Timestamp < startTime {
							break
						}
						ret.Data[interfaceName][metricName] = append(ret.Data[interfaceName][metricName], values[i])
					}
				}
			}
		}
	}

	ret.StatisticsTime = time.Now().Unix()
	return &ret
}

func CleanAllStale(timestamp int64) {
	Locker.Lock()
	defer Locker.Unlock()

	for ip, datas := range globalData.Metrics {
		for metricName, metricData := range datas {
			for interfaceName, values := range metricData.Data {
				i := len(values)
				j := 0
				for ; j < i; j++ {
					if values[j].Timestamp >= timestamp {
						break
					}
				}
				if j >= i {
					j = i - 1
				}
				if j < 0 {
					j = 0
				}
				globalData.Metrics[ip][metricName].Data[interfaceName] = values[j:]
				if len(globalData.Metrics[ip][metricName].Data[interfaceName]) == 0 {
					delete(globalData.Metrics[ip][metricName].Data, interfaceName)
				}
			}
			if len(globalData.Metrics[ip][metricName].Data) == 0 {
				delete(globalData.Metrics[ip], metricName)
			}
		}
		if len(globalData.Metrics[ip]) == 0 {
			delete(globalData.Metrics, ip)
		}
	}
}
