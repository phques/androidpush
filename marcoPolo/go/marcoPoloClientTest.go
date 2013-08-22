// marcoPoloClientTest
package main

import (
	"encoding/json"
	"fmt"
	"net"
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

//----------

func main() {
	// open local UDP port
	fmt.Println("net.DialUDP")

	localUdpAddr := net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}

	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}
	defer udpConn.Close()

	// destination 'marcoPolo' broadcast UDP address
	remoteBroadcastUdpAddr := net.UDPAddr{IP: net.IPv4(255, 255, 255, 255),
		Port: marcoPoloUdpPort}

	// send 'marco'
	fmt.Println("sending marco")

	// create marco.polo msg with empty MarcoPoloMsg json object
	var marcoPoloMsg MarcoPoloMsg
	jsonStr, err := json.MarshalIndent(&marcoPoloMsg, "", "  ")
	marcoMsg := marcoPoloMsgPrefix + string(jsonStr)

	nbBytes, _ := udpConn.WriteToUDP([]byte(marcoMsg), &remoteBroadcastUdpAddr)
	if err != nil {
		fmt.Println("error sending marco ", err)
		return
	}
	_ = nbBytes

}
