package main

import (
	"fmt"
	"net"
	//"beat/utils"
	"flag"
	"encoding/json"
	"log"
	"beat/models"
	"beat/utils"
	"strings"
	"gopkg.in/mgo.v2"
)

const (
	CMD_NewDetector = "NewDetector"
	CMD_HEARTBEAT = "HeartBeat"
	CMD_NewSystemSchedule = "NewSystemSchedule"
	CMD_ScheduleHEARTBEAT = "ScheduleHEARTBEAT"
)

var addr = flag.String("Addr", ":30003", "")

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
	mdb, mSession := utils.GetMgoDbSession()
	defer mSession.Close()
	dta := strings.Split(rmAddr.String(), ":")

	if dat["command"] == CMD_NewDetector {
		device_id := dat["device_id"].(string)
		tp := dat["type"].(string)
		det, e := models.GetDetectorByDeviceId(mdb,device_id, tp)
		if e == mgo.ErrNotFound {
			models.NewDetector(mdb, device_id, tp, dta[0])
		} else {
			det.UpdateByParams(mdb, dta[0])
		}
	}

	if  dat["command"] == CMD_HEARTBEAT {
		device_id := dat["device_id"].(string)
		tp := dat["type"].(string)
		det, e := models.GetDetectorByDeviceId(mdb,device_id,tp)
		if e == mgo.ErrNotFound {
			return
		}
		det.UpdateByStatus(mdb)
	}

	if  dat["command"] == CMD_NewSystemSchedule {
		name := dat["name"].(string)
		_, err := models.GetServiceManage(mdb, name)
		if err == mgo.ErrNotFound {
			models.NewServiceManage(mdb, name)
		}
	}

	if  dat["command"] == CMD_ScheduleHEARTBEAT {
		name := dat["name"].(string)
		manage, err := models.GetServiceManage(mdb, name)
		if err == mgo.ErrNotFound {
			return
		}
		manage.UpdateByStatus(mdb)
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