package flow

import (
	"github.com/SwitchCollector/g"
	"github.com/SwitchCollector/rrdtool"
	"github.com/gaochao1/gosnmp"
	"github.com/gaochao1/sw"
	"log"
	"time"
)

type Flow struct {
	Time            int64   `json:"Time"`
	InFlowQuantity  float64 `json:"InFlowQuantity"`
	OutFlowQuantity float64 `json:"OutFlowQuantity"`
}

type FlowQuantity struct {
	Data []*Flow `json:"Data"`
}

func Search(expire int64) *FlowQuantity {
	inFile := "in.rrd"
	outFile := "out.rrd"
	cf := "AVERAGE"

	endTime := time.Now().Unix()
	startTime := endTime - expire
	step := g.Config().Interval

	var flowQuantity FlowQuantity
	flowQuantity = FlowQuantity{Data: make([]*Flow, 0)}

	inData := rrdtool.FetchFromFile(inFile, cf, startTime, endTime, step)
	outData := rrdtool.FetchFromFile(outFile, cf, startTime, endTime, step)

	if inData == nil || outData == nil {
		return &flowQuantity
	}

	for timestamp := range inData {
		var inValue float64 = 0
		var outValue float64 = 0
		if temp, ok := inData[timestamp]; ok {
			inValue = temp
		}

		if temp, ok := outData[timestamp]; ok {
			outValue = temp
		}
		var flow Flow
		flow = Flow{Time: timestamp, InFlowQuantity: inValue, OutFlowQuantity: outValue}
		flowQuantity.Data = append(flowQuantity.Data, &flow)
	}
	return &flowQuantity
}

func getFlow(ip string, community string, oid string, timeout int) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in get flow", r)
		}
	}()
	method := "get"
	var snmpPDUs []gosnmp.SnmpPDU
	var err error
	for i := 0; i < 3; i++ {
		snmpPDUs, err = sw.RunSnmp(ip, community, oid, method, timeout)
		if len(snmpPDUs) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err == nil {
		for _, pdu := range snmpPDUs {
			return pdu.Value.(uint64), err
		}
	}

	return 0, err
}

func collectAndFlushFlow(ip string, community string, oid string, timeout int, timestamp int64, filename string) {
	if flow, err := getFlow(ip, community, oid, timeout); err == nil {
		item := []*rrdtool.Item{&rrdtool.Item{DsType: "COUNTER", Step: g.Config().Interval, Timestamp: timestamp, Value: float64(flow)}}
		rrdtool.FlushrrdToFile(filename, item)
	}
}

func collectFlowQuantity() {
	log.Println("start collect flow.")
	sh := g.Config().Switch
	timestamp := time.Now().Unix()
	go collectAndFlushFlow(sh.Ip, sh.Community, sh.InFlowOid, sh.Timeout, timestamp, "in.rrd")
	go collectAndFlushFlow(sh.Ip, sh.Community, sh.OutFlowOid, sh.Timeout, timestamp, "out.rrd")
}

func Collect() {
	interval := time.Duration(g.Config().Interval)
	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-ticker.C:
			collectFlowQuantity()
		}
	}
}
