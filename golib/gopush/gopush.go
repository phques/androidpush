// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// package gopush provides the Go functionality for AndroidPush
package gopush

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/phques/mppq"
)

const (
	configFilename string = "config.json"
)

var (
	initDone     bool = false
	started      bool = false
	mppqProvider *mppq.Provider
	config       *Config

	AppFilesDir    string // directory where our app's file are
	ConfigFilepath string // path of our config file (inside appFilesDir)
)

// InitParam holds the info pased from Android app to init gopush
// (dup in goInterface ! because of circular ref))
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
	initAppFilesDir(param.AppFilesDir)

	// Create initial config.json in appFilesDir if does not exists
	if _, err := os.Stat(ConfigFilepath); err != nil {
		if err = createConfigFile(param); err != nil {
			return err
		}
	}

	// load config
	log.Println("loading config from ", ConfigFilepath)
	config = &Config{}
	config.Load(ConfigFilepath)

	//## debug
	log.Printf("config %+v\n", config)

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
	if err := registerMppqService("androidPush"); err != nil {
		return err
	}

	//## test debug
	http.HandleFunc("/androidPush/config", ServeHTTPConfig)

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
func initAppFilesDir(appFilesDir_ string) {

	AppFilesDir = appFilesDir_

	// setup config file path
	ConfigFilepath = filepath.Join(AppFilesDir, configFilename)
	log.Print("config file:", ConfigFilepath)
}

// Create config file from InitParam
func createConfigFile(param *InitParam) error {
	log.Println("creating config file ", ConfigFilepath)
	cfg := createConfigFromInitParam(param)
	return cfg.Save(ConfigFilepath)
}

// createConfigFromInitParam creates a *Config from a *InitParam
func createConfigFromInitParam(param *InitParam) *Config {
	cfg := &Config{}
	cfg.AppFilesDir = param.AppFilesDir
	cfg.Devicename = param.Devicename

	cfg.AddDir("Books", param.Books)
	cfg.AddDir("DCIM", param.DCIM)
	cfg.AddDir("Documents", param.Documents)
	cfg.AddDir("Downloads", param.Downloads)
	cfg.AddDir("Movies", param.Movies)
	cfg.AddDir("Music", param.Music)
	cfg.AddDir("Pictures", param.Pictures)

	return cfg
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

// ServeHTTPConfig handles HTTP GET & PUT for our config file
// curl localhost:1440/androidPush/config
// curl --upload-file ./config.json http://localhost:1440/androidPush/config
func ServeHTTPConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeHTTPConfig", r.Method)

	// GET config file
	if r.Method == "GET" {
		log.Println("ServeHTTPConfig GET ", ConfigFilepath)
		http.ServeFile(w, r, ConfigFilepath)
		return
	}

	// PUT: save config file
	if r.Method == "PUT" {
		log.Println("ServeHTTPConfig PUT ", ConfigFilepath)

		// open output file
		outfile, err := os.Create(ConfigFilepath)
		if err != nil {
			log.Printf("ServeHTTPConfig, error creating file [%v]: %v\n", ConfigFilepath, err)
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer outfile.Close()

		// copy data to output file
		written, err := io.Copy(outfile, r.Body)
		if err != nil {
			log.Printf("ServeHTTPConfig, error copying to file [%v]: %v\n", ConfigFilepath, err)
			return
		}
		log.Printf("ServeHTTPConfig, wrote %v bytes\n", written)
	}

	// invalid method
	log.Println("ServeHTTPConfig, invalid method")
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
