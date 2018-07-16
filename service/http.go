package service

import (
	"encoding/json"
	"github.com/SwitchCollector/g"
	"github.com/SwitchCollector/service/flow"
	"github.com/SwitchCollector/service/visit"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func FlowQuantityBytes(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse data error", http.StatusBadRequest)
		return
	}
	period := strings.TrimSpace(r.Form.Get("Period"))
	expire, err := strconv.ParseInt(period, 10, 64)
	if err != nil {
		http.Error(w, "param Period format error.", http.StatusBadRequest)
		return
	}
	flowQuantity := flow.Query(expire * 60)
	ret, err1 := json.Marshal(flowQuantity)
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(ret)

}

func VisitLog(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse data error", http.StatusBadRequest)
		return
	}
	period := strings.TrimSpace(r.Form.Get("Period"))
	expire, err := strconv.ParseInt(period, 10, 64)
	if err != nil {
		http.Error(w, "param Period format error.", http.StatusBadRequest)
		return
	}
	visitLog := visit.Search(expire * 60)
	ret, err1 := json.Marshal(visitLog)
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(ret)
}

func GetDeviceList(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "parse data error", http.StatusBadRequest)
		return
	}
	deviceList := g.GetDeviceList()
	ret, err := json.Marshal(deviceList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(ret)
}

func StartHttpServer() {
	log.Println("starting http server....")
	http.HandleFunc("/netdevice/devicelist", GetDeviceList)
	http.ListenAndServe(g.Config().Http.Addr, nil)
}
