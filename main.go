package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	srvAddr         = "224.0.0.1:9999"
	maxDatagramSize = 8192
)

var Name string

func main() {
	if len(os.Args) != 2 {
		fmt.Println("example: ./main name")
		return
	}
	Name = os.Args[1]
	go ping(srvAddr)
	serveMulticastUDP(srvAddr, msgHandler)
}

func ping(a string) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	for i := 1; i < 100; i++ {
		c.Write([]byte(Name))
		time.Sleep(5 * time.Second)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(hex.Dump(b[:n]))
}

func serveMulticastUDP(a string, h func(*net.UDPAddr, int, []byte)) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	l.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(src, n, b)
	}
}
