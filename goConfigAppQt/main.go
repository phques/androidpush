package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/phques/androidpush/goConfigAppQt/qmlutil"
	"github.com/phques/mppq"
	"gopkg.in/qml.v1"
)

// WndObjects holds qml objects / 'controls' from our window
// gets filled by qmlutil.FetchObjects
type WndObjects struct {
	ProvidersMdl qml.Object `QML:"providersMdl"` // providers list model
	QueryButton  qml.Object `QML:"queryButton"`
}

// MainWnd is our main window struct
type MainWnd struct {
	objs   WndObjects
	window *qml.Window

	mppqQuery *mppq.Query
	waitGroup sync.WaitGroup // to wait for loopQuery goroutine
	services  map[string]*mppq.ServiceDef
}

// Provider cooresponds to the 'model' for the Providers QML list
type Provider struct {
	Name    string
	Address string
}

//-----

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// createMainWnd creates a useable MainWnd
func createMainWnd() *MainWnd {
	m := &MainWnd{}
	m.services = make(map[string]*mppq.ServiceDef)
	return m
}

func run() error {
	// create new qml engine, load qml file
	engine := qml.NewEngine()
	component, err := engine.LoadFile("goConfigAppQt.qml")
	if err != nil {
		return err
	}

	// create MainWnd and qml window
	mainWnd := createMainWnd()
	mainWnd.window = component.CreateWindow(nil)

	// fill mainWnd.objs with qml.Object from the window ('controls', etc)
	err = qmlutil.FetchObjects(&mainWnd.objs, mainWnd.window.Common)
	if err != nil {
		log.Println("failed to get objects from window:", err)
		return err
	}

	// hook mainWnd.onQueryButtClicked for query button click
	mainWnd.objs.QueryButton.On("clicked", mainWnd.onQueryButtClicked)
	//engine.Context().SetVar("mainWnd", mainWnd)

	mainWnd.window.Show()
	mainWnd.window.Wait()
	return nil
}

// onQueryButtClicked is called when the QueryButton is clicked
func (w *MainWnd) onQueryButtClicked() {
	// do we have a query running ?
	if w.mppqQuery != nil && !w.mppqQuery.Done() {
		w.stopQuery()
	} else {
		w.startQuery()
	}
}

// startQuery starts a new mppq query for androidpush
func (w *MainWnd) startQuery() {
	log.Println("start query")
	// create new mppq query for "androidPush"
	w.mppqQuery = mppq.NewQuery("androidPush", false)

	// start it
	if err := w.mppqQuery.Start(); err != nil {
		//##TODO popup with error
		log.Printf("Error starting mppq query:", err)
		w.mppqQuery = nil
	} else {
		// launch goroutine that reads back found services
		go w.loopQuery()
		w.waitGroup.Add(1)

		//change button to "Stop query"
		w.objs.QueryButton.Set("text", "Stop Query")
	}
}

// stoQuery stops an ongoing mppq query
func (w *MainWnd) stopQuery() {
	log.Println("stop query")
	// stop query
	w.mppqQuery.Stop()

	// wait until goroutine is done
	w.waitGroup.Wait()

	// change button to "Start Query"
	w.objs.QueryButton.Set("text", "Start Query")

	log.Println("stopQuery out (loopQuery has exited)")
}

// loopQuery is a goroutine that receives found services from mppq query
func (w *MainWnd) loopQuery() {
	log.Println("loopQuery in")
	defer log.Println("loopQuery out")
	defer w.waitGroup.Done()
	for {
		select {
		case service, ok := <-w.mppqQuery.ServiceCh:
			if ok {
				log.Println("got service :", service)
				// add to set of services and add to UI list if new
				// use the %v value of service as key
				serviceStr := fmt.Sprintf("%v", service)
				if w.services[serviceStr] == nil {
					// new service, add to map
					w.services[serviceStr] = service
					// add to UI list
					address := fmt.Sprintf("%v:%v", service.RemoteIP, service.HostPort)
					w.AddProvider(Provider{service.ProviderName, address})
					//##TODO: how to match entry in UI with entry in map !!?
				}
			} else {
				//channel closed, stop
				return
			}
		}
	}
}

// AddProvider adds a new Provider to QML providers list
func (w *MainWnd) AddProvider(p Provider) {
	// create json object string, then call providersMdl.myAppend to append to list
	//##(workaround)
	dicStr := fmt.Sprintf(`{"name": "%v", "address": "%v"}`, p.Name, p.Address)
	w.objs.ProvidersMdl.Call("myAppend", dicStr)
}
