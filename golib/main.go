// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/phques/androidpush/golib/goInterface/gen"
	"github.com/phques/androidpush/golib/gopush"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/bind/java"
)

var standalone = flag.Bool("standalone", false, "running as a stand alone pgm")

func main() {
	log.Println("main")
	flag.Parse()
	app.Run(app.Callbacks{Start: start, Stop: stop})
}

func start() {
	log.Println("main.start")
	java.Init()

	if *standalone {
		param := makeStdaloneParam()
		if err := gopush.Init(param); err == nil {
			gopush.Start()
		}
	}
}

func stop() {
	log.Println("main.stop")
	if *standalone {
		//ps: not really required since the app is closing anyways ;-)
		gopush.Stop()
	}
}

// makeStdaloneParam create a InitParam for use when running standalone
// config file in "files" subdir, all root dirs under home/philippe
func makeStdaloneParam() *gopush.InitParam {
	param := &gopush.InitParam{}
	param.Devicename, _ = os.Hostname()
	param.AppFilesDir = "files"
	param.Books = "/home/philippe/Books"
	param.DCIM = "/home/philippe/Pictures"
	param.Documents = "/home/philippe/Documents"
	param.Downloads = "/home/philippe/Downloads"
	param.Movies = "/home/philippe/Movies"
	param.Music = "/home/philippe/Music"
	param.Pictures = "/home/philippe/Pictures"
	return param
}
