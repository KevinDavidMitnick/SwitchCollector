package visit

import (
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
}

type UdpData struct {
	sync.RWMutex
	Data map[string]*AccessIp
}

type VisitLog struct {
	Data           []*AccessIp `json:"Data"`
	StatisticsTime int64       `json:"StatisticsTime"`
}

type IpHistory struct {
	sync.RWMutex
	History map[string][]int64 `json:"history"`
}

var (
	Ips *IpHistory
)

func (ips *IpHistory) save(ip string, timestamp int64) {
	ips.Lock()
	defer ips.Unlock()
	if _, ok := ips.History[ip]; ok {
		ips.History[ip] = append(ips.History[ip], timestamp)
	} else {
		ips.History[ip] = []int64{timestamp}
	}
}

func (ips *IpHistory) search(expire int64) *VisitLog {
	ips.RLock()
	defer ips.RUnlock()
	now := time.Now().Unix()
	startTime := now - expire
	var data []*AccessIp = make([]*AccessIp, 0)

	for ip, history := range ips.History {
		len := len(history)
		var accessIp AccessIp
		accessIp.IP = ip
		accessIp.VisitCount = 0
		for i := len - 1; i >= 0; i-- {
			if startTime <= history[i] {
				accessIp.VisitCount += 1
			}
		}
		if accessIp.VisitCount > 0 {
			data = append(data, &accessIp)
		}
	}

	return &VisitLog{Data: data, StatisticsTime: now}
}

func (ips *IpHistory) cleanStaleData() {
	ips.Lock()
	defer ips.Unlock()
	fmt.Println("start clean stale data.")
	expire := g.Config().Expire
	startTime := time.Now().Unix() - int64(expire)

	for ip, history := range ips.History {
		len := len(history)
		if history[len-1] < startTime {
			delete(ips.History, ip)
			break
		}
		last := len - 1
		for last >= 0 {
			if history[last] < startTime {
				break
			}
			last--
		}
		if last >= 0 {
			ips.History[ip] = history[last:]
		}
	}

}

func (ips *IpHistory) size() int {
	ips.RLock()
	defer ips.RUnlock()

	return len(ips.History)
}

func NewInstance() *UdpData {
	len := g.Config().Expire
	return &UdpData{Data: make(map[string]*AccessIp, len)}
}

func Search(expire int64) *VisitLog {
	return Ips.search(expire)
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
	Ips.save(strings.TrimSpace(ip), timestamp)

	//fmt.Println("recv data is:" + string(data))
}

func NewVisitData() {
	Ips = &IpHistory{
		History: make(map[string][]int64),
	}
}

func CleanStale() {
	expire := time.Duration(g.Config().Expire)
	ticker := time.NewTicker(expire * time.Second)
	for {
		select {
		case <-ticker.C:
			Ips.cleanStaleData()
		}
	}
}
