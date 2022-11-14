package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"net"


//	"github.com/charlesap/sybil/pkg/lodge"
)


func udpResponse(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	time := time.Now().String()
	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))

	udpServer.WriteTo([]byte(responseStr), addr)
}

func slingUDP() {
	udpServer, err := net.ResolveUDPAddr("udp", ":1053")

	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		println("Listen failed:", err.Error())
		os.Exit(1)
	}

	//close the connection
	defer conn.Close()

	_, err = conn.Write([]byte("This is a UDP message"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	println(string(received))
}

func HandleUDP(){

	// listen to incoming udp packets
	udpServer, err := net.ListenPacket("udp", ":1969")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	for {
		buf := make([]byte, 1024)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go udpResponse(udpServer, addr, buf)
	}

}


