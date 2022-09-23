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
	maxDatagramSize = 8192
)

var Name string
var hn string

func main() {
	if len(os.Args) != 3 {
		fmt.Println("example: ./main name [224.0.0.1]:9999")
		return
	}
	Name = os.Args[1]
	srvAddr := os.Args[2]
	hn, _ = os.Hostname()
	go ping(srvAddr)
	serveMulticastUDP(srvAddr, msgHandler)
}

func ping(a string) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	lsnr, err := net.Listen("tcp", ":0")
	c, err := net.DialUDP("udp", nil, addr)
	for i := 1; i < 1000; i++ {
		a, _ := net.InterfaceByName("wlan0")

		if err != nil {
			fmt.Println("Error listening:", err)
			os.Exit(1)
		}
		c.Write([]byte(Name + fmt.Sprint(a.Addrs()) + fmt.Sprint(lsnr.Addr())))
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
	err = l.SetReadBuffer(maxDatagramSize)
	if err != nil {
		return
	}
	for {
		b := make([]byte, maxDatagramSize*4)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		h(src, n, b)
	}
}
