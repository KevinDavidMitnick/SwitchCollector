package device

//Device interface, for example switch,firewall

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/SwitchCollector/core/scheduler"
	"github.com/SwitchCollector/core/store"
	"github.com/SwitchCollector/g"
	"github.com/SwitchCollector/service/funcs"
	log "github.com/sirupsen/logrus"
	"time"
)

type MetricDevice struct {
	Ip           string               `json:"ip"`
	Community    string               `json:"community"`
	Version      string               `json:"version"`
	Class        string               `json:"class"`
	Type         string               `json:"type"`
	Metrics      map[string]*g.Metric `json:"metrics"`
	Infos        map[string]*g.Metric `json:"infos"`
	MultiMetrics map[string]*g.Metric `json:"multimetrics"`
	MultiInfos   map[string]*g.Metric `json:"multiinfos"`
	Timeout      int64                `json:"timeout"`
	Interval     int64                `json:"interval"`
	Uuid         string               `json:"uuid"`
}
type Device struct {
	tasks     []*MetricDevice      `json:"tasks"`
	scheduler *scheduler.Scheduler `json:"scheduler"`
}

type Executer struct {
	scheduler.Object
	Ip         string `json:"ip"`
	Community  string `json:"community"`
	Version    string `json:"version"`
	Oid        string `json:"oid"`
	Interval   int64  `json:"interval"`
	DataType   string `json:"datatype"`
	Timeout    int64  `json:"timeout"`
	Name       string `json:"name"`
	MetricType string `json:"metrictype"`
	MetricName string `json:"metricname"`
	Timestamp  int64  `json:"timestamp"`
	Uuid       string `json:"uuid"`
}

type CmdbResponse struct {
	Message string         `json:"message"`
	Code    int            `json:"code"`
	Data    []*g.NetDevice `json:"data"`
}

func (e *Executer) PingCheck() {
	value, _ := funcs.Ping(e.Ip, e.Timeout)
	e.saveToGD(value)
	if g.Config().Backend.Enabled {
		e.saveToBackend(value)
	}
}

func (e *Executer) PingLatency() {
	_, value := funcs.Ping(e.Ip, e.Timeout)
	e.saveToGD(value)
	if g.Config().Backend.Enabled {
		e.saveToBackend(value)
	}
}

func (e *Executer) saveToBackend(value interface{}) {
	data := make([]map[string]interface{}, 0)

	uuid, ip, metricName, metricType, dataType, timestamp, interval := e.Uuid, e.Ip, e.MetricName, e.MetricType, e.DataType, e.Timestamp, e.Interval
	if metricName == "" {
		metricName = "switch." + e.Name
	}
	switch metricType {
	case "metrics":
		elem := make(map[string]interface{})
		if uuid == "" {
			elem["endpoint"] = ip
		} else {
			elem["endpoint"] = uuid
		}
		elem["metric"] = metricName
		elem["timestamp"] = timestamp
		elem["step"] = interval
		elem["counterType"] = dataType
		elem["tags"] = ""
		elem["value"] = value
		data = append(data, elem)
	case "multimetrics":
		indexNameMap := g.GetIndexNameMap()
		for k, v := range value.(map[string]interface{}) {
			interfaceName := indexNameMap[ip][k]
			elem := make(map[string]interface{})
			if uuid == "" {
				elem["endpoint"] = ip
			} else {
				elem["endpoint"] = uuid
			}
			elem["metric"] = metricName
			elem["timestamp"] = timestamp
			elem["step"] = interval
			elem["counterType"] = dataType
			elem["tags"] = "iface=" + interfaceName
			elem["value"] = v
			data = append(data, elem)
		}
	}
	if len(data) > 0 {
		buf, err := json.Marshal(data)
		if err != nil || len(buf) == 0 {
			log.Println("send json marshal err,or data len is 0 , data is:", data)
			return
		}
		if store.GetStoreStatus() {
			funcs.PushToFalcon(g.Config().Backend.Addr, buf)
		} else {
			store := store.GetStore()
			store.Update(buf)
		}
	}
}

func (e *Executer) saveToGD(value interface{}) {
	gData := g.GetGlobalData()
	indexNameMap := g.GetIndexNameMap()

	g.Locker.Lock()
	defer g.Locker.Unlock()

	ip, name, metricType, dataType, timestamp, interval := e.Ip, e.Name, e.MetricType, e.DataType, e.Timestamp, e.Interval
	if gData.Metrics[ip] == nil {
		gData.Metrics[ip] = make(map[string]*g.MetricData)
	}
	if gData.Metrics[ip][name] == nil {
		gData.Metrics[ip][name] = new(g.MetricData)
	}
	gData.Metrics[ip][name].MetricType = metricType
	if gData.Metrics[ip][name].Data == nil {
		gData.Metrics[ip][name].Data = make(map[string][]*g.DataValue)
	}
	switch metricType {
	case "metrics", "infos":
		if dataType == "GAUGE" {
			dataValue := g.DataValue{LastValue: value, Value: value, Timestamp: timestamp}
			if gData.Metrics[ip][name].Data["liucong"] == nil {
				gData.Metrics[ip][name].Data["liucong"] = []*g.DataValue{&dataValue}
			} else {
				gData.Metrics[ip][name].Data["liucong"] = append(gData.Metrics[ip][name].Data["liucong"], &dataValue)
			}
		} else {
			if gData.Metrics[ip][name].Data["liucong"] == nil {
				dataValue := g.DataValue{LastValue: value, Value: int64(0), Timestamp: timestamp}
				gData.Metrics[ip][name].Data["liucong"] = []*g.DataValue{&dataValue}
			} else {
				len := len(gData.Metrics[ip][name].Data["liucong"])
				curvalue := (value.(int64) - gData.Metrics[ip][name].Data["liucong"][len-1].LastValue.(int64)) / interval
				if curvalue < 0 {
					curvalue = value.(int64) / interval
				}
				dataValue := g.DataValue{LastValue: value, Value: curvalue, Timestamp: timestamp}
				gData.Metrics[ip][name].Data["liucong"] = append(gData.Metrics[ip][name].Data["liucong"], &dataValue)
			}
		}
	case "multimetrics", "multiinfos":
		if dataType == "GAUGE" {
			for k, v := range value.(map[string]interface{}) {
				dataValue := g.DataValue{LastValue: v, Value: v, Timestamp: timestamp}
				if gData.Metrics[ip][name].Data[indexNameMap[ip][k]] == nil {
					gData.Metrics[ip][name].Data[indexNameMap[ip][k]] = []*g.DataValue{&dataValue}
				} else {
					gData.Metrics[ip][name].Data[indexNameMap[ip][k]] = append(gData.Metrics[ip][name].Data[indexNameMap[ip][k]], &dataValue)
				}
			}
		} else {
			for k, v := range value.(map[string]interface{}) {
				if gData.Metrics[ip][name].Data[indexNameMap[ip][k]] == nil {
					dataValue := g.DataValue{LastValue: v, Value: int64(0), Timestamp: timestamp}
					gData.Metrics[ip][name].Data[indexNameMap[ip][k]] = []*g.DataValue{&dataValue}
				} else {
					len := len(gData.Metrics[ip][name].Data[indexNameMap[ip][k]])
					curvalue := (v.(int64) - gData.Metrics[ip][name].Data[indexNameMap[ip][k]][len-1].LastValue.(int64)) / interval
					if curvalue < 0 {
						curvalue = v.(int64) / interval
					}
					dataValue := g.DataValue{LastValue: v, Value: curvalue, Timestamp: timestamp}
					gData.Metrics[ip][name].Data[indexNameMap[ip][k]] = append(gData.Metrics[ip][name].Data[indexNameMap[ip][k]], &dataValue)
				}
			}
		}
	}
}

func (e *Executer) CollectData() {
	querier := funcs.GetQuerier(e.Ip, e.Community, e.Version, e.Timeout)
	defer querier.Close()

	switch e.MetricType {
	case "metrics", "infos":
		value, _ := querier.GetMetricValue(e.Oid)
		e.saveToGD(value)
		if g.Config().Backend.Enabled {
			e.saveToBackend(value)
		}
	case "multimetrics", "multiinfos":
		value, _ := querier.GetBulkMetricValue(e.Oid)
		e.saveToGD(value)
		if g.Config().Backend.Enabled {
			e.saveToBackend(value)
		}
	default:
		log.Println("not right metric type:", e.MetricType)
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

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buf).Decode(dst)
}

func mergeMetrics(dev *g.NetDevice, metricT *g.MetricTemplate) *MetricDevice {
	var device MetricDevice
	device.Ip = dev.Ip
	device.Community = dev.Community
	device.Version = dev.Version
	device.Class = dev.Class
	device.Type = dev.Type
	device.Uuid = dev.Uuid

	metrics := make(map[string]*g.Metric)
	infos := make(map[string]*g.Metric)
	multiMetrics := make(map[string]*g.Metric)
	multiInfos := make(map[string]*g.Metric)
	if metricT != nil {
		deepCopy(&metrics, metricT.Metrics)
		deepCopy(&infos, metricT.Infos)
		deepCopy(&multiMetrics, metricT.MultiMetrics)
		deepCopy(&multiInfos, metricT.MultiInfos)
	}

	if dev.Extension != nil && dev.Extension.Enabled {
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

func buildIndexNameMap(ip string, community string, version string, timeout int64) {
	indexNameMap := g.GetIndexNameMap()

	g.Locker.Lock()
	defer g.Locker.Unlock()
	if indexNameMap[ip] == nil {
		indexNameMap[ip] = make(map[string]string)
	}
	querier := funcs.GetQuerier(ip, community, version, timeout)
	defer querier.Close()
	tempMap, err := querier.GetBulkMetricValue(".1.3.6.1.2.1.2.2.1.2")
	if err != nil {
		log.Println(err.Error())
	}
	for key, value := range tempMap {
		indexNameMap[ip][key] = value.(string)
	}
	//log.Println("map is:", indexNameMap)
}

func (device *Device) InitTasks() {
	metricM := g.MetricT()
	netDevs := g.NetDevs()
	for ip := range netDevs {
		dev := netDevs[ip]
		metricDevice := mergeMetrics(dev, metricM[dev.Type])
		device.tasks = append(device.tasks, metricDevice)
		buildIndexNameMap(ip, dev.Community, dev.Version, metricM[dev.Type].Timeout)
	}
}

func (device *Device) InitScheduler() {
	device.scheduler.Lock()
	defer device.scheduler.Unlock()
	for _, metricDevice := range device.tasks {
		metricsL := map[string]map[string]*g.Metric{"metrics": metricDevice.Metrics,
			"infos": metricDevice.Infos, "multimetrics": metricDevice.MultiMetrics, "multiinfos": metricDevice.MultiInfos}
		for metricType, metrics := range metricsL {
			for name, metric := range metrics {
				var executer Executer
				executer.Ip = metricDevice.Ip
				executer.Community = metricDevice.Community
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
				executer.MetricName = metric.MetricName
				executer.Uuid = metricDevice.Uuid

				device.scheduler.Queue[executer.Interval] = append(device.scheduler.Queue[executer.Interval], &executer)
			}
		}
	}
}

func (device *Device) Collect() {
	device.scheduler.Scheduler()
}

func (device *Device) CleanStale() {
	expire := time.Duration(g.Config().Expire)
	ticker := time.NewTicker(expire * time.Second)
	for {
		select {
		case <-ticker.C:
			timestamp := time.Now().Unix() - int64(g.Config().Expire)
			log.Println("start clean flow stale data")
			g.CleanAllStale(timestamp)
			log.Println("finish start clean flow stale data")
		}
	}
}

func (device *Device) StopAll() {
	device.scheduler.Lock()
	defer device.scheduler.Unlock()

	for key := range device.scheduler.Queue {
		device.scheduler.Queue[key] = nil
	}
}

func (device *Device) Diff(devices []*g.NetDevice) (incExecuter []*Executer, decExecuter []*Executer) {
	device.scheduler.RLock()
	defer device.scheduler.RUnlock()
	for interval, objects := range device.scheduler.Queue {
		for i, object := range objects {
			var flag bool = true
			for _, dev := range devices {
				if object.(*Executer).Uuid == dev.Uuid {
					flag = false
					device.scheduler.Queue[interval][i].(*Executer).Ip = dev.Ip
					device.scheduler.Queue[interval][i].(*Executer).Community = dev.Community
					device.scheduler.Queue[interval][i].(*Executer).Version = dev.Version
				}
			}
			if flag {
				decExecuter = append(decExecuter, object.(*Executer))
			}
		}
	}
	for _, dev := range devices {
		var flag bool = true
		for interval, objects := range device.scheduler.Queue {
			for i, object := range objects {
				if object.(*Executer).Uuid == dev.Uuid {
					flag = false
					device.scheduler.Queue[interval][i].(*Executer).Ip = dev.Ip
					device.scheduler.Queue[interval][i].(*Executer).Community = dev.Community
					device.scheduler.Queue[interval][i].(*Executer).Version = dev.Version
				}
			}
		}
		if flag {
			metricM := g.MetricT()
			if metricM[dev.Type] == nil {
				break
			}
			metricDevice := mergeMetrics(dev, metricM[dev.Type])
			metricsL := map[string]map[string]*g.Metric{"metrics": metricDevice.Metrics,
				"infos": metricDevice.Infos, "multimetrics": metricDevice.MultiMetrics, "multiinfos": metricDevice.MultiInfos}
			for metricType, metrics := range metricsL {
				for name, metric := range metrics {
					var executer Executer
					executer.Ip = metricDevice.Ip
					executer.Community = metricDevice.Community
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
					executer.MetricName = metric.MetricName
					executer.Uuid = metricDevice.Uuid
					incExecuter = append(incExecuter, &executer)
				}
			}
		}
	}
	return
}

func (device *Device) Increase(incExecuter []*Executer) {
	device.scheduler.Lock()
	defer device.scheduler.Unlock()

	indexNameMap := g.GetIndexNameMap()
	var newIntervals []int64
	for _, executer := range incExecuter {
		if indexNameMap[executer.Ip] == nil {
			buildIndexNameMap(executer.Ip, executer.Community, executer.Version, executer.Timeout)
		}
		interval := executer.Interval
		if len(device.scheduler.Queue[interval]) == 0 {
			newIntervals = append(newIntervals, interval)
		}
		device.scheduler.Queue[interval] = append(device.scheduler.Queue[interval], executer)
	}
	for _, val := range newIntervals {
		go device.scheduler.Run(val)
	}
}

func (device *Device) Decrease(decExecuter []*Executer) {
	device.scheduler.Lock()
	defer device.scheduler.Unlock()
	uuids := make(map[int64][]string)
	for _, executer := range decExecuter {
		uuids[executer.Interval] = append(uuids[executer.Interval], executer.Uuid)
	}
	for interval := range uuids {
		var objects []scheduler.Object
		for _, object := range device.scheduler.Queue[interval] {
			var flag bool = true
			for _, uuid := range uuids[interval] {
				if object.(*Executer).Uuid == uuid {
					flag = false
				}
			}
			if flag {
				objects = append(objects, object)
			}
		}
		device.scheduler.Queue[interval] = objects
	}
}

func (device *Device) UpdateScheduler() {
	if !g.Config().Cmdb.Enabled {
		return
	}

	interval := time.Duration(g.Config().Interval)
	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println("start update scheduler data")
			var response CmdbResponse
			cmdbUrl := g.Config().Cmdb.Addr
			data, err := funcs.GetData(cmdbUrl)
			if err != nil {
				continue
			}
			err1 := json.NewDecoder(bytes.NewBuffer(data)).Decode(&response)
			if err1 != nil {
				continue
			}
			if response.Code != 0 || response.Message != "ok" {
				continue
			}
			if len(response.Data) == 0 {
				device.StopAll()
				continue
			}
			var devices []*g.NetDevice
			for _, dev := range response.Data {
				if dev.Ip == "" || dev.Community == "" || dev.Class == "" ||
					dev.Type == "" || dev.Uuid == "" {
					continue
				}
				devices = append(devices, dev)
			}
			incExecuter, decExecuter := device.Diff(devices)
			device.Increase(incExecuter)
			device.Decrease(decExecuter)
			log.Println("finish update scheduler data")
		}
	}
}

func (device *Device) UpdateStoreStatus() {
	_, err := funcs.GetData(g.Config().Backend.Check)
	if err == nil {
		store.UpdateStoreStatus(true)
	} else {
		store.UpdateStoreStatus(false)
	}
}

func (device *Device) FlushStore() {
	if !g.Config().Backend.Enabled {
		return
	}
	interval := time.Duration(g.Config().Interval)
	s := store.GetStore()
	defer s.Close()
	for {
		timestamp := time.Now().Unix() - int64(g.Config().Expire)
		data := make([]map[string]interface{}, 0)
		queue := make(chan []byte, g.Config().Interval)
		s.CleanStale(timestamp, data)
		s = store.GetStore()
		device.UpdateStoreStatus()

		go func() {
			for {
				data := <-queue
				if data == nil {
					break
				}
				funcs.PushToFalcon(g.Config().Backend.Addr, data)
			}
		}()
		for data := s.Read(); store.GetStoreStatus() && data != nil; data = s.Read() {
			queue <- data
		}
		close(queue)
		time.Sleep(interval * time.Second)
	}
}
