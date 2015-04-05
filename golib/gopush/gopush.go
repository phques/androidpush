// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// package gopush provides the Go functionality for AndroidPush
package gopush

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/phques/mppq"
)

const (
	mppqServiceName string = "androidPush"
)

var (
	initDone     bool = false
	started      bool = false
	mppqProvider *mppq.Provider
	config       *Config

	appFilesDir    string // directory where our app's file are
	configFilename string // path of our config file (inside appFilesDir)
)

// InitParam holds the info pased from Android app to init gopush
// (dup in goInterface ! because of circular ref))
type InitParam struct {
	Devicename  string // reported name to mppq service query responses
	AppFilesDir string // app's files dir, we store config file there

	// config file directories, used to populate config file 1st time
	Books     string // can be empty, will try as sibling of Documents
	Camera    string // can be empty, will try under DCIM/Camera
	DCIM      string
	Documents string
	Downloads string
	Movies    string
	Music     string
	Pictures  string
}

//-----

// Init initializes the Gopush lib
func Init(param *InitParam) error {
	log.Println("gopush.Init")
	// already done ?
	if initDone {
		log.Println(" .. already done")
		return nil
	}
	initDone = true

	// set app's files dir & config filename
	if err := initAppFilesDir(param.AppFilesDir); err != nil {
		return err
	}

	// Create initial config.json in appFilesDir if does not exists
	// (we checked that file's dir is Ok in Init)
	if !isFile(configFilename) {
		if err := createConfigFile(param); err != nil {
			return err
		}
	}

	// load config
	log.Println("loading config from ", configFilename)
	config = &Config{}
	if err := config.Load(configFilename); err != nil {
		log.Println("config.Load error : ", err)
		return err
	}

	return nil
}

// Start() starts http & mppq servers, registers androidPush service with mppq.
func Start() error {
	log.Println("gopush.Start")
	if started {
		log.Println(" .. already done")
		return nil
	}
	started = true

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
	if err := registerMppqService(mppqServiceName); err != nil {
		return err
	}

	// register GET/PUT handler for config file
	http.HandleFunc("/androidPush/config", serveHTTPConfig)

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

// initAppFilesDir initializes the app's files dir values
func initAppFilesDir(appFilesDir_ string) error {
	_, err := os.Stat(appFilesDir_)
	if err != nil {
		log.Println("Can't find app files dir: ", err)
		return err
	}

	appFilesDir = appFilesDir_

	// setup config file path
	configFilename = filepath.Join(appFilesDir, "config.json")
	log.Print("config file: ", configFilename)

	return nil
}

// Create config file from InitParam
func createConfigFile(param *InitParam) error {
	log.Println("creating config file ", configFilename)
	cfg := createConfigFromInitParam(param)
	return cfg.Save(configFilename)
}

// createConfigFromInitParam creates a *Config from a *InitParam
func createConfigFromInitParam(param *InitParam) *Config {
	cfg := &Config{}
	cfg.AppFilesDir = param.AppFilesDir
	cfg.Devicename = param.Devicename

	checkBooksDir(param)
	checkCameraDir(param)

	cfg.AddDir("Books", param.Books)
	cfg.AddDir("Camera", param.Camera)
	cfg.AddDir("DCIM", param.DCIM)

	cfg.AddDir("Documents", param.Documents)
	cfg.AddDir("Downloads", param.Downloads)
	cfg.AddDir("Movies", param.Movies)
	cfg.AddDir("Music", param.Music)
	cfg.AddDir("Pictures", param.Pictures)

	return cfg
}

// checkBooksDir tries to find Books dir at same level as
// Documents dir in param if it is empty
func checkBooksDir(param *InitParam) {
	// if Books dir empty
	if len(param.Books) == 0 && len(param.Documents) != 0 {
		// try "Books" at same level as "Documents"
		booksDir, found := lookForSiblingDir(param.Documents, "Books")
		if found {
			// ok found .../Books, save in param
			param.Books = booksDir
		}

		if len(param.Books) != 0 {
			log.Println("init, found books dir :", param.Books)
		} else {
			log.Println("init, no books dir found")
		}
	}
}

// checkCameraDir tries to find Camera dir under the DCIM dir
func checkCameraDir(param *InitParam) {
	// if Camera dir empty
	if len(param.Camera) == 0 {
		// try "Camera" under DCIM
		cameraDir := filepath.Join(param.DCIM, "Camera")
		if isDir(cameraDir) {
			// ok found
			param.Camera = cameraDir
		}

		if len(param.Camera) != 0 {
			log.Println("init, found Camera dir :", param.Camera)
		} else {
			log.Println("init, no Camera dir found")
		}
	}
}

func lookForSiblingDir(dir, lookFor string) (foundDir string, found bool) {
	parentDir := filepath.Dir(dir)
	siblingDir := filepath.Join(parentDir, lookFor)
	if isDir(siblingDir) {
		return siblingDir, true
	}
	return "", false
}

func isDir(dir string) bool {
	fileinfo, err := os.Stat(dir)
	if err == nil {
		if fileinfo.IsDir() {
			return true
		}
	} else if !os.IsNotExist(err) { // ??
		log.Println(err)
	}
	return false
}

func isFile(filename string) bool {
	fileinfo, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			// other Stat() error
			log.Printf("Init, error: %v\n", configFilename, err)
			return false
		}
		return false
	}
	if fileinfo.Mode()&os.ModeType != 0 {
		// special or Dir
		return false
	}
	return true
}

// register a service we provide with mppq
func registerMppqService(serviceName string) error {

	log.Println("registerMppqService", serviceName)

	// register a service (mppqProvider must be started)
	deviceName := config.Devicename
	if len(deviceName) == 0 {
		deviceName, _ = os.Hostname() // returns 'localhost' on my Nexus 5/7
	}
	err := mppqProvider.AddService(mppq.ServiceDef{
		ServiceName:  serviceName,
		ProviderName: deviceName,
		HostPort:     httpListenPort,
		Protocol:     "jsonhttp",
	})

	return err
}
