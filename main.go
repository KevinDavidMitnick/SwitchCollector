package main

import (
	"flag"
	"fmt"
	"github.com/SwitchCollector/g"
	"github.com/SwitchCollector/service"
	"os"
)

func init() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.LoadTemplatesConfig()
	g.LoadNetDevices()
}

func main() {
	service.Init()
	go service.CollectFlow()
	if g.Config().Http.Enabled {
		go service.StartHttpServer()
	}
	if g.Config().Udp.Enabled {
		go service.StartUdpServ()
	}

	select {}
}
