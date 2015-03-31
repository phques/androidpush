// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/phques/androidpush/golib/goInterface/gen"
	"github.com/phques/androidpush/golib/gopush"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/bind/java"
)

var standalone = flag.Bool("standalone", false, "running as a stand alone pgm")

func main() {
	fmt.Println("main")
	flag.Parse()
	app.Run(app.Callbacks{Start: start, Stop: stop})
}

func start() {
	fmt.Println("main.start")
	java.Init()

	if *standalone {
		param := gopush.InitParam{}
		param.Hostname, _ = os.Hostname()
		param.AppFilesDir = "files"
		gopush.Init(&param)
		gopush.Start()
	}
}

func stop() {
	fmt.Println("main.stop")
	if *standalone {
		//ps: not really required since the app is closing anyways ;-)
		gopush.Stop()
	}
}
