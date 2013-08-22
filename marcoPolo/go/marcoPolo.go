// udpStr.go
package main

import (
	"fmt"
	"net"
	//"time"
	"encoding/json"
	"strings"
)

type Version struct {
	Major int
	Minor int
}

type MarcoPoloMsg struct {
	Version    Version
	Action     string
	Name       string
	OptPayload string // optional broadcast msg payload (json string)
}

//---------

var marcoPoloUdpPort int = 4444
var marcoPoloMsgPrefix = "marco.polo:"

//---------

func isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

func main() {

	// open local UDP port
	fmt.Println("net.DialUDP")

	localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: marcoPoloUdpPort}

	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}
	defer udpConn.Close()

	// wait for msg
	data := make([]byte, 1024*16)
	read, remoteAddr, err := udpConn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("error reading udp socket ", err)
		return
	}
	udpMarcoPolo := data[:read]

	fmt.Printf("read %d bytes '%s'\n", read, udpMarcoPolo)
	fmt.Printf("from udp (remote %s)\n", remoteAddr.String())

	// check for marco polo msg
	if strings.HasPrefix(string(udpMarcoPolo), marcoPoloMsgPrefix) {
		// marco.polo:{JSON MarcoPoloMsg}
		marcoPoloMsgJson := udpMarcoPolo[len(marcoPoloMsgPrefix):]
		fmt.Println("marcoPoloMsgJson:", string(marcoPoloMsgJson))

		// unmarshall json string to MarcoPoloMsg
		var marcoPoloMsg MarcoPoloMsg
		err = json.Unmarshal(marcoPoloMsgJson, &marcoPoloMsg)
		if err == nil {
			fmt.Printf("marcoPoloMsg: %+v\n", marcoPoloMsg)
		} else {
			fmt.Println(err)
			return
		}
	}
}
