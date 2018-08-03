package service

import (
	"github.com/SwitchCollector/g"
	"github.com/SwitchCollector/service/visit"
	log "github.com/sirupsen/logrus"
	"net"
)

func StartUdpServ() {
	// create upd socket
	log.Println("starting udp server....")
	udpAddr, err := net.ResolveUDPAddr("udp4", g.Config().Udp.Addr)
	if err != nil {
		log.Fatalln(err)
	}
	socket, err1 := net.ListenUDP("udp4", udpAddr)
	if err1 != nil {
		log.Fatalln(err)
	}
	defer socket.Close()

	visit.Init()
	go visit.CleanStale()
	for {
		// read data.
		data := make([]byte, 1024)
		size, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			log.Println("read data from udp failed,remote ip :" + string(remoteAddr.IP))
			continue
		}
		go visit.HandleUdpData(data, size)
	}
}
