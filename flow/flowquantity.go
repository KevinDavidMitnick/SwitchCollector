package flow

import (
	"github.com/SwitchCollector/g"
	"github.com/yubo/rrdlite"
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

func FetchFromFile(filename string, cf string, start, end int64, step int) map[int64]float64 {
	var data map[int64]float64
	startT := time.Unix(start, 0)
	endT := time.Unix(end, 0)
	stepT := time.Duration(step) * time.Second

	fetchRes, err := rrdlite.Fetch(filename, cf, startT, endT, stepT)
	if err != nil {
		return data
	}

	defer fetchRes.FreeValues()

	values := fetchRes.Values()
	size := len(values)
	data = make(map[int64]float64, size)

	startTs := fetchRes.Start.Unix()
	stepS := fetchRes.Step.Seconds()

	for i, val := range values {
		ts := startTs + int64(i+1)*int64(stepS)
		data[ts] = float64(val)
	}

	return data
}

func Search(expire int64) *FlowQuantity {
	inFile := "in.rrd"
	outFile := "out.rrd"
	cf := "GAUGE"

	endTime := time.Now().Unix()
	startTime := endTime - expire
	step := g.Config().Interval

	flowQuantity := &FlowQuantity{Data: make([]*Flow, 0)}

	inData := FetchFromFile(inFile, cf, startTime, endTime, step)
	outData := FetchFromFile(outFile, cf, startTime, endTime, step)

	if inData == nil || outData == nil {
		return flowQuantity
	}

	for timestamp := range inData {
		flow := &Flow{Time: timestamp, InFlowQuantity: inData[timestamp], OutFlowQuantity: outData[timestamp]}
		flowQuantity.Data = append(flowQuantity.Data, flow)
	}
	return flowQuantity
}
