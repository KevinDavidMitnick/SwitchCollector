package device

//Device interface, for example switch,firewall

import (
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
	Interval     int                  `json:"interval"`
}
type Devices struct {
	task []*MetricDevice `json:"task"`
}

var (
	devices Devices
)

func mergeMetrics(dev *g.NetDevice, metricT *g.MetricTemplate) *MetricDevice {
	var device MetricDevice
	device.Ip = dev.Ip
	device.Community = dev.Community
	device.Port = dev.Port
	device.Version = dev.Version
	device.Class = dev.Class
	device.Type = dev.Type

	metrics := metricT.Metrics
	infos := metricT.Infos
	multiMetrics := metricT.MultiMetrics
	multiInfos := metricT.MultiInfos
	if dev.Extension.Enabled {
		for key, value := range dev.Extension.Metrics {
			if metrics == nil {
				metrics = make(map[string]*g.Metric)
			}
			metrics[key] = value
		}
		for key, value := range dev.Extension.Infos {
			if infos == nil {
				infos = make(map[string]*g.Metric)
			}
			infos[key] = value
		}
		for key, value := range dev.Extension.MultiMetrics {
			if multiMetrics == nil {
				multiMetrics = make(map[string]*g.Metric)
			}
			multiMetrics[key] = value
		}
		for key, value := range dev.Extension.MultiInfos {
			if multiInfos == nil {
				multiInfos = make(map[string]*g.Metric)
			}
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

func Init() {
	if devices.task == nil {
		devices.task = make([]*MetricDevice, 0)
	}
	metricM := g.MetricT()
	netDevs := g.NetDevs()
	for ip := range netDevs {
		dev := netDevs[ip]
		device := mergeMetrics(dev, metricM[dev.Type])
		devices.task = append(devices.task, device)
	}
}

func (d *Devices) Update() {

}

func (d *Devices) Collect() {

}

func (d *Devices) Flush() {

}

func (d *Devices) CleanStale() {

}
