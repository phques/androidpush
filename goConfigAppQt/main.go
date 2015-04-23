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
	ProvidersMdl  qml.Object `QML:"providersMdl"` // providers list model
	QueryButton   qml.Object `QML:"queryButton"`
	MessageDialog qml.Object `QML:"messageDialog"`
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

// newMainWnd creates a useable MainWnd
func newMainWnd() *MainWnd {
	m := &MainWnd{}
	m.services = make(map[string]*mppq.ServiceDef)
	return m
}

// setupGUI creates the main window, fetches QML objects etc
func (w *MainWnd) setupGUI(component qml.Object) error {
	w.window = component.CreateWindow(nil)

	// fill w.objs with qml.Object from the window ('controls', etc)
	err := qmlutil.FetchObjects(&w.objs, w.window.Common)
	if err != nil {
		log.Println("failed to get objects from window:", err)
		return err
	}

	// hook w.onQueryButtClicked for query button click
	w.objs.QueryButton.On("clicked", w.onQueryButtClicked)

	//engine.Context().SetVar("mainWnd", w)

	return nil

}
func run() error {
	// create new qml engine, load qml file
	engine := qml.NewEngine()
	component, err := engine.LoadFile("goConfigAppQt.qml")
	if err != nil {
		return err
	}

	// create MainWnd & GUI
	// create MainWnd and qml window
	mainWnd := newMainWnd()
	if err := mainWnd.setupGUI(component); err != nil {
		return err
	}

	mainWnd.window.Show()
	mainWnd.window.Wait()
	return nil
}

// ShowMsg open a QML MessageDialog with title & text set
func (w *MainWnd) ShowMsg(title, text string) {
	box := w.objs.MessageDialog
	box.Set("title", title)
	box.Set("text", text)
	box.Call("open")
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
		log.Printf("Error starting mppq query:", err)
		w.ShowMsg("Start Query", "Error :"+err.Error())
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

	// loop waiting for received service definitions
	for {
		select {
		case service, ok := <-w.mppqQuery.ServiceCh:
			if !ok {
				//channel closed, stop
				return
			}

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
		}
	}
}

// AddProvider adds a new Provider to QML providers list
func (w *MainWnd) AddProvider(p Provider) {
	//##(workaround)
	// create json object string, then call providersMdl.myAppend to append to list
	dicStr := fmt.Sprintf(`{"name": "%v", "address": "%v"}`, p.Name, p.Address)
	w.objs.ProvidersMdl.Call("myAppend", dicStr)
}
