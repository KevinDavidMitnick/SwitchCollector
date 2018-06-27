package visit

import (
	"encoding/json"
	"fmt"
	"github.com/SwitchCollector/g"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type AccessIp struct {
	IP         string `json:"IP"`
	VisitCount int    `json:"VisitCount"`
	LastTime   int64  `json:"LastTime"`
}

type UdpData struct {
	sync.RWMutex
	Data map[string]*AccessIp
}

type VisitLog struct {
	Data           []*AccessIp `json:"Data"`
	StatisticsTime int64       `json:"StatisticsTime"`
}

var VisitData *UdpData

func (udp *UdpData) save(ip string, timestamp int64) {
	udp.Lock()
	defer udp.Unlock()
	if _, ok := udp.Data[ip]; ok {
		udp.Data[ip].VisitCount += 1
		udp.Data[ip].LastTime = timestamp
	} else {
		udp.Data[ip] = &AccessIp{IP: ip, VisitCount: 1, LastTime: timestamp}
	}
}

func (udp *UdpData) get(ip string) *AccessIp {
	udp.RLock()
	defer udp.RUnlock()
	if data, ok := udp.Data["ip"]; ok {
		return data
	}
	return nil
}

func (udp *UdpData) display() {
	udp.RLock()
	defer udp.RUnlock()

	if data, err := json.Marshal(udp.Data); err == nil {
		fmt.Println(string(data))
	}
}

func (udp *UdpData) search(expire int64) *VisitLog {
	udp.RLock()
	defer udp.RUnlock()
	now := time.Now().Unix()
	startTime := now - expire
	var data []*AccessIp = make([]*AccessIp, 0)

	udp.display()

	for _, accessIp := range udp.Data {
		if accessIp.LastTime >= startTime {
			data = append(data, accessIp)
		}
	}

	return &VisitLog{Data: data, StatisticsTime: now}
}

func (udp *UdpData) cleanStaleData() {
	udp.Lock()
	defer udp.Unlock()
	fmt.Println("start clean stale data.")
	expire := g.Config().Expire
	startTime := time.Now().Unix() - int64(expire)

	for ip, accessIp := range udp.Data {
		if accessIp.LastTime < startTime {
			delete(udp.Data, ip)
		}
	}

}

func (udp *UdpData) size() int {
	udp.RLock()
	defer udp.RUnlock()

	return len(udp.Data)
}

func NewInstance() *UdpData {
	len := g.Config().Expire
	return &UdpData{Data: make(map[string]*AccessIp, len)}
}

func Search(expire int64) *VisitLog {
	return VisitData.search(expire)
}

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

	for {
		// read data.
		data := make([]byte, 1024)
		size, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			log.Println("read data from udp failed,remote ip :" + string(remoteAddr.IP))
			continue
		}
		go handleUdpData(data, size)
	}
}

func handleUdpData(data []byte, size int) {
	timestamp := time.Now().Unix()
	str := string(data)
	if strings.Contains(str, "src_addr") == false {
		log.Println("err package format:" + str)
		return
	}
	packet := strings.Split(str, ";")
	ip_port := packet[7]
	if strings.Contains(ip_port, "src_addr") == false {
		log.Println("err package format:" + ip_port)
		return
	}
	ip := strings.Split(ip_port, ":")[1]
	VisitData.save(strings.TrimSpace(ip), timestamp)

	//fmt.Println("recv data is:" + string(data))
	fmt.Println("visit data size is:", VisitData.size())
	//VisitData.display()
}

func NewVisitData() {
	VisitData = NewInstance()
}

func CleanStale() {
	expire := time.Duration(g.Config().Expire)
	ticker := time.NewTicker(expire * time.Second)
	for {
		select {
		case <-ticker.C:
			VisitData.cleanStaleData()
		}
	}
}
