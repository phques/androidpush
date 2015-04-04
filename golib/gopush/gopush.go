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
	configFilename  string = "config.json"
	mppqServiceName string = "androidPush"
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
	if err := initAppFilesDir(param.AppFilesDir); err != nil {
		return err
	}

	// Create initial config.json in appFilesDir if does not exists
	// (we checked that file's dir is Ok in Init)
	_, err := os.Stat(ConfigFilepath)
	if err != nil {
		// if doesnt exist, create it
		if os.IsNotExist(err) {
			if err = createConfigFile(param); err != nil {
				return err
			}
		} else {
			// other stat() error
			log.Printf("stat(%v) error: %v\n", ConfigFilepath, err)
			return err
		}
	}

	// load config
	log.Println("loading config from ", ConfigFilepath)
	config = &Config{}
	if err = config.Load(ConfigFilepath); err != nil {
		log.Println("config.Load error : ", err)
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
	http.HandleFunc("/androidPush/config", ServeHTTPConfig)

	//## test debug
	docsDir := config.Dirs["Documents"][0]
	log.Println(docsDir)
	docServer := http.FileServer(http.Dir(docsDir))
	//docServer := http.FileServer(http.Dir("/home/philippe/Documents"))
	http.Handle("/androidPush/Documents/", http.StripPrefix("/androidPush/Documents/", docServer))

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
		log.Println(err)
		return err
	}

	AppFilesDir = appFilesDir_

	// setup config file path
	ConfigFilepath = filepath.Join(AppFilesDir, configFilename)
	log.Print("config file:", ConfigFilepath)

	return nil
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
	deviceName := config.Devicename
	if deviceName == "" {
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

// ServeHTTPConfig handles HTTP GET & PUT for our config file
// GET: curl localhost:1440/androidPush/config -o config.json
// PUT: curl --upload-file ./config.json http://localhost:1440/androidPush/config
func ServeHTTPConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeHTTPConfig", r.Method)

	// GET config file
	if r.Method == "GET" {
		log.Println("  GET ", ConfigFilepath)
		http.ServeFile(w, r, ConfigFilepath)
		return
	}

	// PUT: save config file
	if r.Method == "PUT" {
		log.Println("  PUT ", ConfigFilepath)

		if saveConfig(w, r) {
			// re-read config (will that work, outfile is not closed?)
			config.Load(ConfigFilepath)

			// re-register mppq androidPush (devicename might have changed)
			registerMppqService(mppqServiceName)
		}

		return
	}

	// invalid method
	log.Println("ServeHTTPConfig, invalid method", r.Method)
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func saveConfig(w http.ResponseWriter, r *http.Request) bool {
	// open output file
	outfile, err := os.Create(ConfigFilepath)
	if err != nil {
		log.Printf("ServeHTTPConfig, error creating file [%v]: %v\n", ConfigFilepath, err)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return false
	}
	defer outfile.Close()

	// copy data to output file
	written, err := io.Copy(outfile, r.Body)
	if err != nil {
		log.Printf("ServeHTTPConfig, error copying to file [%v]: %v\n", ConfigFilepath, err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return false
	}

	log.Printf("ServeHTTPConfig, wrote %v bytes\n", written)
	return true
}
