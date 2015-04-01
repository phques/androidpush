// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package gopush

import (
	"log"
	"net"
	"net/http"
	"strconv"
)

var (
	httpListener   net.Listener
	httpListenPort int
)

// StartHTTP starts a http.Serve() go routine, listening on a sys allocated local port
// The listener and it's port are saved in httpListener / httpListenPort
func StartHTTP() error {
	// listen, on system assigned port
	ln, err := net.Listen("tcp", ":1440") //## debug
	//ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Printf("Failed to start HTTP Server : %v\n", err)
		return err
	}

	// save listener port
	addr := ln.Addr()
	_, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		httpListener.Close()
		log.Printf("Failed to start HTTP Server, error getting listener port : %v\n", err)
		return err
	}
	portNo, err := strconv.Atoi(port)
	if err != nil {
		httpListener.Close()
		log.Printf("Failed to start HTTP Server, error getting listener port : %v\n", err)
	}

	//## debug
	log.Println("http server port", portNo)

	// save listener & port
	httpListener = ln
	httpListenPort = portNo

	// start serving
	go http.Serve(ln, nil)
	return nil
}
