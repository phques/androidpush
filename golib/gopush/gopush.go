// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// package gopush provides the Go functionality for AndroidPush
package gopush

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/phques/mppq"
	"golang.org/x/mobile/app"
)

const (
	configFilename string = "config.json"
)

var (
	InitDone     bool = false
	mppqProvider *mppq.Provider

	AppFilesDir    string // directory where our app's file are
	ConfigFilepath string // path of our config file (inside appFilesDir)
)

// InitParam holds the info pased from Android app to init gopush
// (dup in goInterface ! because of circular ref))
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

//-----

// Init initializes the Gopush lib
func Init(param *InitParam) error {
	log.Println(*param)

	if err := initAppFilesDir(param.AppFilesDir); err != nil {
		return err
	}

	return nil
}

// Start() starts http & mppq servers, registers androidPush service with mppq.
func Start() error {
	log.Println("gopush.Start")

	// start http server
	if err := StartHTTP(); err != nil {
		return err
	}

	// create/start mppq provider
	mppqProvider = mppq.NewProvider()
	if err := mppqProvider.Start(); err != nil {
		return err
	}

	// register androidPush
	if err := registerMppqService("androidPush"); err != nil {
		return err
	}

	return nil
}

// Stop stops the mppq provider
func Stop() error {
	log.Println("gopush.Stop")

	if mppqProvider == nil {
		return errors.New("Stop, Provider is not running")
	}

	return mppqProvider.Stop()
}

// InitAppFilesDir initializes the app's files dir & copies config file there the 1st time
func initAppFilesDir(appFilesDir_ string) error {
	// already done ?
	if InitDone {
		return nil
	}
	InitDone = true

	//## debug
	dir, _ := os.Getwd()
	log.Printf("cwd: %v\n", dir)

	AppFilesDir = appFilesDir_

	// setup config file path
	ConfigFilepath = filepath.Join(AppFilesDir, configFilename)
	log.Print("config file:", ConfigFilepath)

	// create initial (copy from assets) config.json in appFilesDir if does not exists
	// does config file exist in app files dir?
	if _, err := os.Stat(ConfigFilepath); err != nil {
		return copyConfigFile()
	}

	return nil
}

//--- utils -----

// copy config file from assets to app filesdir
func copyConfigFile() (err error) {
	// open src config file from assets
	srcFile, err := app.Open(configFilename)
	if err != nil {
		log.Printf("copyConfigFile, error opening source : %v\n", err)
		return
	}
	defer srcFile.Close()

	// create/open dest config file
	destFile, err := os.Create(ConfigFilepath)
	if err != nil {
		log.Printf("copyConfigFile, error opening dest : %v\n", err)
		return
	}
	defer destFile.Close()

	// copy
	nbCopied, err := io.Copy(destFile, srcFile)
	if err == nil {
		log.Printf("copyConfigFile, copied %v bytes\n", nbCopied)
	} else {
		log.Printf("copyConfigFile, error copying : %v\n", err)
	}

	return nil
}

// register a service we provide with mppq
func registerMppqService(serviceName string) error {

	log.Println("registerMppqService", serviceName)

	// register a service (mppqProvider must be started)
	//## PQ TODO use 'deviceName' from config
	providerName, _ := os.Hostname() // returns 'localhost' on my Nexus 5/7
	err := mppqProvider.AddService(mppq.ServiceDef{
		ServiceName:  serviceName,
		ProviderName: providerName,
		HostPort:     httpListenPort,
		Protocol:     "jsonhttp",
	})

	return err
}
