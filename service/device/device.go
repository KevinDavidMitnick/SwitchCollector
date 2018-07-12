package device

//Device interface, for example switch,firewall

import (
	"fmt"
	"github.com/SwitchCollector/core/scheduler"
	"github.com/SwitchCollector/g"
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
	Ip        string `json:"ip"`
	Community string `json:"community"`
	Port      int    `json:"port"`
	Version   string `json:"version"`
	Oid       string `json:"oid"`
	Interval  int64  `json:"interval"`
	DataType  string `json:"datatype"`
	Timeout   int    `json:"timeout"`
	Name      string `json:"name"`
}

func (e *Executer) Run() {
	fmt.Println("run:", e.Ip, e.Oid)

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

func (device *Device) InitTasks() {
	metricM := g.MetricT()
	netDevs := g.NetDevs()
	for ip := range netDevs {
		dev := netDevs[ip]
		metricDevice := mergeMetrics(dev, metricM[dev.Type])
		device.tasks = append(device.tasks, metricDevice)
	}
}

func (device *Device) Update() {

}

func (device *Device) InitScheduler() {
	for _, metricDevice := range device.tasks {
		metricsL := []map[string]*g.Metric{metricDevice.Metrics,
			metricDevice.Infos, metricDevice.MultiMetrics, metricDevice.MultiInfos}
		for _, metrics := range metricsL {
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
