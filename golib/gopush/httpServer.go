// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package gopush

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
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
	ln, err := net.Listen("tcp", ":0")
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

// used to intercept ResponseWriter.Write(buff []byte)
// when GETting a file so we can keep track of progress
type MyWriter struct {
	http.ResponseWriter
	filename string
	size     int64
	written  int64
}

func NewMyWriter(w http.ResponseWriter, filename string) *MyWriter {
	mw := &MyWriter{ResponseWriter: w, filename: filename}

	st, err := os.Stat(mw.filename)
	if err == nil {
		mw.size = st.Size()
	} else {
		log.Printf("NewMyWriter, failed to get size of [%v]", mw.filename)
	}

	return mw
}

func (mw *MyWriter) Write(buff []byte) (int, error) {
	buffLen := int64(len(buff))
	mw.written += buffLen
	log.Printf("MyWriter.Write, %v/%v bytes (%v) \n", mw.written, mw.size, buffLen)
	return mw.ResponseWriter.Write(buff)
}

// serveHTTPConfig handles HTTP GET & PUT for our config file
// GET: curl localhost:1440/androidPush/config -o config.json
// PUT: curl --upload-file ./config.json http://localhost:1440/androidPush/config
func serveHTTPConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeHTTPConfig", r.Method)

	// GET config file
	if r.Method == "GET" {
		log.Println("  GET ", configFilename)
		mw := NewMyWriter(w, configFilename)
		http.ServeFile(mw, r, configFilename)
		//		http.ServeFile(w, r, ConfigFilepath)
		return
	}

	// PUT: save config file
	if r.Method == "PUT" {
		log.Println("  PUT ", configFilename)

		if saveConfig(w, r) {
			// re-read config
			//##PQ: config could be in use somewhere else !
			//##todo: notify user to restart app ?
			/*
				if err := config.Load(ConfigFilepath); err != nil {
					log.Printf("error reloading config: %v", err)
				} else {
					// re-register mppq androidPush (devicename might have changed)
					registerMppqService(mppqServiceName)
				}
			*/
		}
		return
	}

	// invalid method
	log.Println("ServeHTTPConfig, invalid method", r.Method)
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// saveConfig saves our config file on a PUT
func saveConfig(w http.ResponseWriter, r *http.Request) bool {
	// open output file
	outfile, err := os.Create(configFilename)
	if err != nil {
		log.Printf("ServeHTTPConfig, error creating file [%v]: %v\n", configFilename, err)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return false
	}
	defer outfile.Close()

	// copy data to output file
	written, err := io.Copy(outfile, r.Body)
	if err != nil {
		log.Printf("ServeHTTPConfig, error copying to file [%v]: %v\n", configFilename, err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return false
	}

	log.Printf("ServeHTTPConfig, wrote %v bytes\n", written)
	return true
}
