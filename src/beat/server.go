package main

import (
	"fmt"
	"net"
	//"beat/utils"
	"flag"
	"encoding/json"
	"log"
)

const (
	CMD_NewDetector = "NewDetector"
	CMD_HEARTBEAT = "HeartBeat"
)

var addr = flag.String("Addr", ":30001", "")

func init() {
	flag.Parse()
}

func main() {
	//Resolving address
	udpAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		fmt.Println("Resolving address Error: ", err)
	}

	// Build listining connections
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Build listining connections Error: ", err)
	}
	defer conn.Close()
	//recvBuff := make([]byte, 1024)
	//ackBuf := []byte("ok!")
	for {
		log.Println("Ready to receive packets!")
		recvUDPMsg(conn)
	}
}

func recvUDPMsg(conn *net.UDPConn){
	recvBuff := make([]byte, 1024)
	ackBuf := []byte("ok!")
	rn, rmAddr, err := conn.ReadFromUDP(recvBuff)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	res := string(recvBuff[:rn])
	byt := []byte(res)
	var dat map[string]interface{}
	if err = json.Unmarshal(byt, &dat); err != nil {
		ackBuf = []byte("NO!")
		fmt.Println("已经报错了")
		fmt.Println(err)
		return
	}

	fmt.Println(dat["command"])
	fmt.Println(dat["type"])
	fmt.Println(dat["device_id"])

	fmt.Printf("<<< Packet received from: %s, data: %s\n", rmAddr.String(), string(recvBuff[:rn]))

	_, err = conn.WriteToUDP(ackBuf, rmAddr)
	if err != nil {
		log.Println("Error:", err)
		return
	}
}