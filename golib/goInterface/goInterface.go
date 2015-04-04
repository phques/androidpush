// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// Package goInterface is the Go API for the Android Java app
package goInterface

import (
	"log"

	"github.com/phques/androidpush/golib/gopush"
)

// InitParam holds the info pased from Android app to init gopush
// (dup in gopush ! because of circular ref))
type InitParam struct {
	Devicename  string // reported name to mppq service query responses
	AppFilesDir string // app's files dir, we store config file there

	// config file directories, used to populate config file 1st time
	Books     string
	DCIM      string // for the Camera
	Documents string
	Downloads string
	Movies    string
	Music     string
	Pictures  string
}

//------

// NewInitParam creates an empty InitParam (required in Java)
func NewInitParam() *InitParam {
	return new(InitParam)
}

// Init initializes the Gopush library
func Init(param *InitParam) error {
	log.Println("goInterface.Init")
	return gopush.Init(param.dupInitParam())
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

// dupInitParam copies InitParam into gopush.InitParam
func (i *InitParam) dupInitParam() *gopush.InitParam {
	p := new(gopush.InitParam)
	p.Devicename = i.Devicename
	p.AppFilesDir = i.AppFilesDir
	p.Music = i.Music
	p.Downloads = i.Downloads
	p.Documents = i.Documents
	p.Pictures = i.Pictures
	p.Movies = i.Movies
	p.Books = i.Books
	p.DCIM = i.DCIM
	return p
}
