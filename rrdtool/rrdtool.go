package rrdtool

import (
	"errors"
	"github.com/toolkits/file"
	"github.com/yubo/rrdlite"
	"math"
	"strconv"
	"time"
)

// RRA.Point.Size
const (
	RRA1PointCnt = 4320 // 1m一个点存12h *2 *3
)

type Item struct {
	DsType    string  `json:"dstype"`
	Step      int     `json:"step"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

func create(filename string, item *Item) error {
	now := time.Now()
	start := now.Add(time.Duration(-24) * time.Hour)
	step := uint(item.Step)

	c := rrdlite.NewCreator(filename, start, step)
	c.DS("metric", item.DsType, step*2, 'U', 'U')

	// 设置各种归档策略
	// 默认1分钟一个点存 3d
	c.RRA("AVERAGE", 0, 1, RRA1PointCnt)
	return c.Create(true)
}

func FlushrrdToFile(filename string, items []*Item) error {
	if items == nil || len(items) == 0 {
		return errors.New("empty items")
	}

	if !file.IsExist(filename) {
		baseDir := file.Dir(filename)

		err := file.InsureDir(baseDir)
		if err != nil {
			return err
		}

		err = create(filename, items[0])
		if err != nil {
			return err
		}
	}

	return update(filename, items)
}

func update(filename string, items []*Item) error {
	u := rrdlite.NewUpdater(filename)

	for _, item := range items {
		v := math.Abs(item.Value)
		if v > 1e+300 || (v < 1e-300 && v > 0) {
			continue
		}
		if item.DsType == "DERIVE" || item.DsType == "COUNTER" {
			u.Cache(item.Timestamp, int(item.Value))
		} else {
			u.Cache(item.Timestamp, item.Value)
		}
	}

	return u.Update()
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
		if math.IsNaN(val) {
			data[ts] = 0
		} else {
			strval := strconv.FormatFloat(val, 'f', 2, 64)
			data[ts], _ = strconv.ParseFloat(strval, 64)
		}
	}

	return data
}
