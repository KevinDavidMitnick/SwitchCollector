package device

//Device interface, for example switch,firewall

import (
	"github.com/SwitchCollector/g"
)

type MetricDevice struct {
	Ip           string               `json:"ip"`
	Community    string               `json:"community"`
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

func mergeMetrics(dev *g.NetDevice, metricT *MetricTemplate) *MetricDevice {
	var device MetricDevice
	device.Ip = dev.Ip
	device.Community = dev.Community
	device.Class = dev.Class
	device.Type = dev.Type

	metrics := metricT.Metrics
	infos := metricT.Infos
	multiMetrics := metricT.MultiMetrics
	multiInfos := metricT.MultiInfos
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

func (d *Device) Init() {
	if devices.task == nil {
		devices.task = make([]*MetricDevice)
	}
	metricM := g.MetricT()
	netDevs := g.NetDevs()
	for ip := range netDevs {
		dev := netDevs[ip]
		device := mergeMetrics(dev, metricM[dev.Type])
		devices.task = append(devices.task, device)
	}
}

func (d *Device) Update() {

}

func (d *Device) Collect() {

}

func (d *Device) Flush() {

}

func (d *Device) CleanStale() {

}
