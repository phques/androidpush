package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/gonutz/goiup/iup"
	"github.com/gonutz/goiup/iuputil"
	"github.com/phques/mppq"
)

// Will hold Handles to controls,
// filled by iuputil.FetchControls
type MyControls struct {
	MainDialog      *iup.Handle `IUP:"mainDialog"`
	Providers       *iup.Handle `IUP:"providers"`
	QueryButt       *iup.Handle `IUP:"queryButton"`
	StopQueryButt   *iup.Handle `IUP:"stopQueryButton"`
	ProviderDetails *iup.Handle `IUP:"providerDetails"`
}

type App struct {
	MyControls
	query   *mppq.Query
	cmdChan chan func() // channel to send GUI commands to execute
}

var app App

//---------

// Idle callback
// Called from goroutine to execute commands that change the GUI,
// since the GUI stuff must run in the original thread
func idleFunc1() int {

	select {
	case cmd := <-app.cmdChan:
		//		log.Println("got something to do in idle")
		// execute the command we recvd
		cmd()

	case <-time.After(time.Duration(100 * time.Millisecond)):
		// timeout, nothing to do
	}

	return iup.DEFAULT
}

//----------

func createDialog() {
	// load GUI definitions from file
	if errStr := iup.Load("androidGUI.led"); errStr != "" {
		panic(errStr)
	}

	// get controls handles into myControls
	if err := iuputil.FetchControls(&app.MyControls); err != nil {
		panic(err)
		return
	}

	// prepare a channel for the idle callback msgs,
	// start a goroutine to send a msg on the channel after some time
	app.cmdChan = make(chan func())

	// hook our idle func
	iup.SetIdleFunc(idleFunc1)

	// hook controls callbacks
	app.QueryButt.SetCallback("ACTION", app.queryBtnCB)
	app.StopQueryButt.SetCallback("ACTION", app.stopQueryBtnCB)

}

// enableQueryButts enables/disables the query/stop query buttons
func (app *App) enableQueryButts() {
	if app.query == nil {
		app.QueryButt.SetAttribute("ACTIVE", "YES")
		app.StopQueryButt.SetAttribute("ACTIVE", "NO")
	} else {
		app.QueryButt.SetAttribute("ACTIVE", "NO")
		app.StopQueryButt.SetAttribute("ACTIVE", "YES")
	}
}

// queryBtnCB is the callbackfor the Query button
func (app *App) queryBtnCB() int {
	log.Println("queryBtnCB, start query loop")

	// Create a mppq query
	app.query = mppq.NewQuery("androidPush", false)
	app.enableQueryButts()

	// start the query
	err := app.query.Start()
	if err != nil {
		//oops, error, cleanup
		msg := fmt.Sprintf("Failed to start mppq query loop:\n%v", err)
		iup.Message("AndroidPush", msg)

		app.query = nil
		app.enableQueryButts()
	}

	// start goroutine to receive query responses
	go app.loopRecvMppq()

	return iup.DEFAULT
}

// stopQueryBtnCB is the callback for the Stop Query button
func (app *App) stopQueryBtnCB() int {
	log.Println("stopQueryBtnCB, stop query loop")

	app.query.Stop()
	app.query = nil
	app.enableQueryButts()

	return iup.DEFAULT
}

// loopRecvMppq loops reading app.query.ServiceCh
func (app *App) loopRecvMppq() {
	for {
		service, ok := <-app.query.ServiceCh
		if ok {
			serviceStr := fmt.Sprintf("%v", service)
			log.Println(serviceStr)

			// send a command to execute: add item to providers list
			app.cmdChan <- func() {
				app.MyControls.Providers.SetAttribute("APPENDITEM", serviceStr)
			}
		} else {
			log.Println("loopRecvMppq quit")
			break
		}
	}
}

//----------

func main() {
	runtime.LockOSThread()

	iup.Open()
	defer iup.Close()

	createDialog()

	// show dialog and loop until last window closed
	app.MainDialog.Show()

	iup.MainLoop()

	//## test debug
	time.Sleep(time.Millisecond * 250)
}
