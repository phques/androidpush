// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// marcoPoloSrv.go
package main

import (
	"fmt"
	//"net"
	//"time"
	"encoding/json"
	"google.code/p/marcoPoloGo/marcopolo"
	"strings"
)

func main() {

	// open local UDP port
	fmt.Println("open udp connection")

	srvConn, err := marcopolo.OpenServerConn()
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}
	defer srvConn.Close()

	// wait for marco.polo msg

	data := make([]byte, 1024*16)
	read, remoteAddr, err := srvConn.UdpConn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("error reading udp socket ", err)
		return
	}
	udpMarcoPolo := data[:read]

	fmt.Printf("read %d bytes '%s'\n", read, udpMarcoPolo)
	fmt.Printf("from udp (remote %s)\n", remoteAddr.String())

	// check for marco polo msg
	if strings.HasPrefix(string(udpMarcoPolo), marcopolo.MsgPrefix) {
		// get the JSON string (a MarcoPoloMsg)
		// "marco.polo:{JSON MarcoPoloMsg}"
		marcoPoloMsgJson := udpMarcoPolo[len(marcopolo.MsgPrefix):]
		fmt.Println("marcoPoloMsgJson:", string(marcoPoloMsgJson))

		// unmarshall json string to MarcoPoloMsg
		var marcoPoloMsg marcopolo.Msg
		err = json.Unmarshal(marcoPoloMsgJson, &marcoPoloMsg)
		if err == nil {
			fmt.Printf("marcoPoloMsg: %+v\n", marcoPoloMsg)
		} else {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("recv data that was not a marco.polo msg")
	}
}
