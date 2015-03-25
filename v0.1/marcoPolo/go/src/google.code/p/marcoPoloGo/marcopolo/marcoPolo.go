// AndroidPush project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// marcoPolo.go
package marcopolo

import (
	//"fmt"
	"net"
	//"time"
	//"encoding/json"
	//"strings"
)

type Version struct {
	Major int
	Minor int
}

type Msg struct {
	Version    Version
	Action     string
	Name       string
	OptPayload string // optional broadcast msg payload (json string)
}

//---------

// ## temp public
type UdpConn struct {
	UdpConn *net.UDPConn
	UdpAddr net.UDPAddr
}

type ServerConn UdpConn
type ClientConn UdpConn

//---------

const (
	UdpPort   int = 4444 // fixed / hard-coded port for now ;-p
	MsgPrefix     = "marco.polo:"
)

var (
	localAnyAddr = net.IPv4(0, 0, 0, 0)
)

//---------

func OpenServerConn() (serverConn ServerConn, err error) {
	serverConn = ServerConn{}

	// server local udp connection, system assigned address & port = marcoPoloUdpPort
	serverConn.UdpAddr = net.UDPAddr{IP: localAnyAddr, Port: UdpPort}

	// open local udp connection
	serverConn.UdpConn, err = net.ListenUDP("udp4", &serverConn.UdpAddr)
	return
}

func OpenClientConn() (clientConn ClientConn, err error) {
	clientConn = ClientConn{}

	// client local udp socket, system assigned addr & port
	clientConn.UdpAddr = net.UDPAddr{IP: localAnyAddr, Port: 0}

	// open local udp connection
	clientConn.UdpConn, err = net.ListenUDP("udp4", &clientConn.UdpAddr)
	return
}

func (conn *UdpConn) Close() {
	if conn.UdpConn != nil {
		conn.UdpConn.Close()
	}
}

func (conn *ServerConn) Close() {
	(*UdpConn)(conn).Close()
}
