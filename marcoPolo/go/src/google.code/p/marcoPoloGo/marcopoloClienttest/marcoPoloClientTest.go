// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

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

	// local udp socket
	localAnyAddr := net.IPv4(0, 0, 0, 0)
	localUdpAddr := net.UDPAddr{IP: localAnyAddr, Port: 0}
	// destination 'marcoPolo' broadcast UDP address on marcoPoloUdpPort
	broadcastAddr := net.IPv4(255, 255, 255, 255)
	remoteBroadcastUdpAddr := net.UDPAddr{IP: broadcastAddr, Port: marcoPoloUdpPort}

	// open local udp connection
	udpConn, err := net.ListenUDP("udp4", &localUdpAddr)
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}
	defer udpConn.Close()


	// send marcoPolo
	fmt.Println("sending marco")

	// create marco.polo msg with empty MarcoPoloMsg json object
	var marcoPoloMsg MarcoPoloMsg
	jsonStr, _ := json.MarshalIndent(&marcoPoloMsg, "", "  ")
	marcoMsg := marcoPoloMsgPrefix + string(jsonStr)

	// send msg
	nbBytes, err := udpConn.WriteToUDP([]byte(marcoMsg), &remoteBroadcastUdpAddr)
	if err != nil {
		fmt.Println("error sending marco ", err)
		return
	}
	_ = nbBytes

}
