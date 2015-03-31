// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// Package goInterface is the Go API for the Android Java app
package goInterface

import (
	"log"

	"github.com/phques/androidpush/golib/gopush"
)

// circular ref prob between gopush and goInterface !!
// this needs to be here for gobind,
// but I need to pass an InitParam to gopush ... !!??
type InitParam struct {
	Hostname    string // reported name to mppq service query responses
	AppFilesDir string // app's files dir, we store config file there

	// config file directories, used to populate config file 1st time
	Music     string
	Downloads string
	Documents string
	Pictures  string
	Movies    string
	Books     string
	DCIM      string // for the Camera
}

//------

// NewInitParam creates an empty InitParam (required in Java)
func NewInitParam() *InitParam {
	return new(InitParam)
}

// Init initializes the Gopush library
func Init(param *InitParam) {
	log.Println(*param)

	//gopush.Init(param)
}

// Start() starts http & mppq servers, registers androidPush service with mppq.
// NB: InitAppFilesDir should be called 1st
func Start() error {
	log.Println("goInterface.Start")
	return gopush.Start()
}

// Stop stops the mppq provider
func Stop() error {
	log.Println("goInterface.Stop")
	return gopush.Stop()
}
