package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln(err)
	}
	defer c.Close()

	for {
		mt, message, err1 := c.ReadMessage()
		if err1 != nil {
			fmt.Println("read from socket err:", err1.Error())
			break
		}
		go writeMessage(c, mt, message)
	}
}

func writeMessage(conn *websocket.Conn, mt int, msg []byte) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			conn.WriteMessage(mt, msg)
		}
	}
}

func startUdpServ() {
	// create upd socket
	log.Println("starting udp server....")
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 514,
	})
	if err != nil {
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
	fmt.Println(data)
}

func main() {
	fmt.Println("starting server....")
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	go http.ListenAndServe("0.0.0.0:8080", nil)

	go startUdpServ()
	select {}
}
