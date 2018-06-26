package main

import (
	"flag"
	"fmt"
	"github.com/SwitchCollector/g"
	"github.com/SwitchCollector/visit"
	"net/http"
	"os"
)

type Flow struct {
	Time            int64 `json:"Time"`
	InFlowQuantity  int   `json:"InFlowQuantity"`
	OutFlowQuantity int   `json:"OutFlowQuantity"`
}

type FlowQuantity struct {
	Data []*Flow `json:"Data"`
}

func FlowQuantityBytes(w http.ResponseWriter, r *http.Request) {

}

func VisitLog(w http.ResponseWriter, r *http.Request) {

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
	select {}
}
