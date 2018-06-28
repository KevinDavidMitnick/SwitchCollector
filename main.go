package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SwitchCollector/flow"
	"github.com/SwitchCollector/g"
	"github.com/SwitchCollector/visit"
	"net/http"
	"os"
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

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	fmt.Println("starting http server....")
	http.HandleFunc("/FlowQuantityBytes", FlowQuantityBytes)
	http.HandleFunc("/VisitLog", VisitLog)
	go http.ListenAndServe(g.Config().Http.Addr, nil)

	visit.NewVisitData()
	go visit.StartUdpServ()
	go visit.CleanStale()

	go flow.Collect()
	go flow.CleanStale()

	select {}
}
