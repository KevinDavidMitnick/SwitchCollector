package device

//Device interface, for example switch,firewall

import (
	"fmt"
	"github.com/SwitchCollector/core/scheduler"
	"github.com/SwitchCollector/g"
	"github.com/SwitchCollector/service/funcs"
)

type MetricDevice struct {
	Ip           string               `json:"ip"`
	Community    string               `json:"community"`
	Port         int                  `json:"port"`
	Version      string               `json:"version"`
	Class        string               `json:"class"`
	Type         string               `json:"type"`
	Metrics      map[string]*g.Metric `json:"metrics"`
	Infos        map[string]*g.Metric `json:"infos"`
	MultiMetrics map[string]*g.Metric `json:"multimetrics"`
	MultiInfos   map[string]*g.Metric `json:"multiinfos"`
	Timeout      int                  `json:"timeout"`
	Interval     int64                `json:"interval"`
}
type Device struct {
	tasks     []*MetricDevice      `json:"tasks"`
	scheduler *scheduler.Scheduler `json:"scheduler"`
}

type Executer struct {
	scheduler.Object
	Ip         string `json:"ip"`
	Community  string `json:"community"`
	Port       int    `json:"port"`
	Version    string `json:"version"`
	Oid        string `json:"oid"`
	Interval   int64  `json:"interval"`
	DataType   string `json:"datatype"`
	Timeout    int    `json:"timeout"`
	Name       string `json:"name"`
	MetricType string `json:"metrictype"`
	Timestamp  int64  `json:"timestamp"`
}

func (e *Executer) PingCheck() {
	gData := g.GetGlobalData()
	value, _ := funcs.Ping(e.Ip, e.Timeout)
	saveToGD(e.Ip, e.Name, e.Timeout, e.MetricType, e.DataType, e.Timestamp, e.Interval, value)
	fmt.Println("check is:", gData.Metrics[e.Ip][e.Name].Data["liucong"][0].Value.(int64))
}

func (e *Executer) PingLatency() {
	gData := g.GetGlobalData()
	_, value := funcs.Ping(e.Ip, e.Timeout)
	saveToGD(e.Ip, e.Name, e.Timeout, e.MetricType, e.DataType, e.Timestamp, e.Interval, value)
	fmt.Println("latency is:", gData.Metrics[e.Ip][e.Name].Data["liucong"][0].Value.(int64))
}

func saveToGD(ip string, name string, timeout int, metricType string, dataType string, timestamp int64, interval int64, value interface{}) {
	gData := g.GetGlobalData()
	indexNameMap := g.GetIndexNameMap()

	if gData.Metrics[ip] == nil {
		gData.Metrics[ip] = make(map[string]*g.MetricData)
	}
	if gData.Metrics[ip][name] == nil {
		gData.Metrics[ip][name] = new(g.MetricData)
	}
	gData.Metrics[ip][name].MetricType = metricType
	if gData.Metrics[ip][name].Data == nil {
		gData.Metrics[ip][name].Data = make(map[string][]*g.DataValue)
		switch metricType {
		case "metrics", "infos":
			if dataType == "GAUGE" {
				dataValue := g.DataValue{LastValue: value, Value: value, Timestamp: timestamp}
				gData.Metrics[ip][name].Data["liucong"] = []*g.DataValue{&dataValue}
			} else {
				dataValue := g.DataValue{LastValue: value, Value: int64(0), Timestamp: timestamp}
				gData.Metrics[ip][name].Data["liucong"] = []*g.DataValue{&dataValue}
			}
		case "multimetrics", "multiinfos":
			if dataType == "GAUGE" {
				for k, v := range value.(map[string]interface{}) {
					dataValue := g.DataValue{LastValue: v, Value: v, Timestamp: timestamp}
					gData.Metrics[ip][name].Data[indexNameMap[k]] = []*g.DataValue{&dataValue}
				}
			} else {
				for k, v := range value.(map[string]interface{}) {
					dataValue := g.DataValue{LastValue: v, Value: int64(0), Timestamp: timestamp}
					gData.Metrics[ip][name].Data[indexNameMap[k]] = []*g.DataValue{&dataValue}
				}
			}
		}
	} else {
		switch metricType {
		case "metrics", "infos":
			if dataType == "GAUGE" {
				dataValue := g.DataValue{LastValue: value, Value: value, Timestamp: timestamp}
				gData.Metrics[ip][name].Data["liucong"] = append(gData.Metrics[ip][name].Data["liucong"], &dataValue)
			} else {
				len := len(gData.Metrics[ip][name].Data["liucong"])
				curvalue := (value.(int64) - gData.Metrics[ip][name].Data["liucong"][len-1].LastValue.(int64)) / interval
				dataValue := g.DataValue{LastValue: value, Value: curvalue, Timestamp: timestamp}
				gData.Metrics[ip][name].Data["liucong"] = append(gData.Metrics[ip][name].Data["liucong"], &dataValue)
			}
		case "multimetrics", "multiinfos":
			if dataType == "GAUGE" {
				for k, v := range value.(map[string]interface{}) {
					dataValue := g.DataValue{LastValue: v, Value: v, Timestamp: timestamp}
					gData.Metrics[ip][name].Data[indexNameMap[k]] = append(gData.Metrics[ip][name].Data[indexNameMap[k]], &dataValue)
				}
			} else {
				for k, v := range value.(map[string]interface{}) {
					len := len(gData.Metrics[ip][name].Data[indexNameMap[k]])
					curvalue := (v.(int64) - gData.Metrics[ip][name].Data[indexNameMap[k]][len-1].LastValue.(int64)) / interval
					dataValue := g.DataValue{LastValue: v, Value: curvalue, Timestamp: timestamp}
					gData.Metrics[ip][name].Data[indexNameMap[k]] = append(gData.Metrics[ip][name].Data[indexNameMap[k]], &dataValue)
				}
			}
		}
	}
}

func (e *Executer) CollectData() {
	querier := funcs.GetQuerier(e.Ip, e.Port, e.Community, e.Version, e.Timeout)
	defer querier.Close()

	switch e.MetricType {
	case "metrics", "infos":
		value, _ := querier.GetMetricValue(e.Oid)
		saveToGD(e.Ip, e.Name, e.Timeout, e.MetricType, e.DataType, e.Timestamp, e.Interval, value)
	case "multimetrics", "multiinfos":
		value, _ := querier.GetBulkMetricValue(e.Oid)
		saveToGD(e.Ip, e.Name, e.Timeout, e.MetricType, e.DataType, e.Timestamp, e.Interval, value)
	default:
		fmt.Println("not right metric type:", e.MetricType)
	}
}

func (e *Executer) Run(timestamp int64) {
	e.Timestamp = timestamp
	switch e.Oid {
	case "ping_check":
		e.PingCheck()
	case "ping_latency":
		e.PingLatency()
	default:
		e.CollectData()
	}
}

func mergeMetrics(dev *g.NetDevice, metricT *g.MetricTemplate) *MetricDevice {
	var device MetricDevice
	device.Ip = dev.Ip
	device.Community = dev.Community
	device.Port = dev.Port
	device.Version = dev.Version
	device.Class = dev.Class
	device.Type = dev.Type

	metrics := metricT.Metrics
	if metrics == nil {
		metrics = make(map[string]*g.Metric)
	}

	infos := metricT.Infos
	if infos == nil {
		infos = make(map[string]*g.Metric)
	}

	multiMetrics := metricT.MultiMetrics
	if multiMetrics == nil {
		multiMetrics = make(map[string]*g.Metric)
	}

	multiInfos := metricT.MultiInfos
	if multiInfos == nil {
		multiInfos = make(map[string]*g.Metric)
	}

	if dev.Extension.Enabled {
		for key, value := range dev.Extension.Metrics {
			metrics[key] = value
		}
		for key, value := range dev.Extension.Infos {
			infos[key] = value
		}
		for key, value := range dev.Extension.MultiMetrics {
			multiMetrics[key] = value
		}
		for key, value := range dev.Extension.MultiInfos {
			multiInfos[key] = value
		}
	}
	device.Metrics = metrics
	device.Infos = infos
	device.MultiMetrics = multiMetrics
	device.MultiInfos = multiInfos
	device.Timeout = metricT.Timeout
	device.Interval = metricT.Interval
	return &device
}

func GetDevice() *Device {
	var device Device
	if device.tasks == nil {
		device.tasks = make([]*MetricDevice, 0)
	}
	device.scheduler = scheduler.GetScheduler()
	return &device
}

func buildIndexNameMap(ip string, port int, community string, version string, timeout int) {
	indexNameMap := g.GetIndexNameMap()
	if len(indexNameMap) != 0 {
		return
	}
	querier := funcs.GetQuerier(ip, port, community, version, timeout)
	defer querier.Close()
	tempMap, err := querier.GetBulkMetricValue(".1.3.6.1.2.1.2.2.1.2")
	if err != nil {
		fmt.Println(err.Error())
	}
	for key, value := range tempMap {
		indexNameMap[key] = value.(string)
	}
	fmt.Println("map is:", indexNameMap)
}

func (device *Device) InitTasks() {
	metricM := g.MetricT()
	netDevs := g.NetDevs()
	for ip := range netDevs {
		dev := netDevs[ip]
		metricDevice := mergeMetrics(dev, metricM[dev.Type])
		device.tasks = append(device.tasks, metricDevice)
		buildIndexNameMap(ip, dev.Port, dev.Community, dev.Version, metricM[dev.Type].Timeout)
	}
}

func (device *Device) Update() {

}

func (device *Device) InitScheduler() {
	for _, metricDevice := range device.tasks {
		metricsL := map[string]map[string]*g.Metric{"metrics": metricDevice.Metrics,
			"infos": metricDevice.Infos, "multimetrics": metricDevice.MultiMetrics, "multiinfos": metricDevice.MultiInfos}
		for metricType, metrics := range metricsL {
			for name, metric := range metrics {
				var executer Executer
				executer.Ip = metricDevice.Ip
				executer.Community = metricDevice.Community
				executer.Port = metricDevice.Port
				executer.Version = metricDevice.Version
				executer.Interval = metricDevice.Interval
				executer.Oid = metric.Oid
				if metric.Interval > 0 {
					executer.Interval = metric.Interval
				}
				executer.DataType = metric.DataType
				executer.Timeout = metricDevice.Timeout
				executer.Name = name
				executer.MetricType = metricType

				device.scheduler.Queue[executer.Interval] = append(device.scheduler.Queue[executer.Interval], &executer)
			}
		}
	}
}

func (device *Device) Collect() {
	device.scheduler.Scheduler()
	select {}
}

func (device *Device) Flush() {

}

func (device *Device) CleanStale() {

}
