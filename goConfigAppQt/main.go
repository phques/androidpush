package main

import (
	"fmt"
	"github.com/phques/androidpush/goConfigAppQt/qmlutil"
	"os"

	"gopkg.in/qml.v1"
)

// WndObjects holds qml objects / 'controls' from our window
// filled with qmlutil.FetchObjects
type WndObjects struct {
	ProvidersMdl qml.Object `QML:"providersMdl"` // providers list model
	QueryButton  qml.Object `QML:"queryButton"`
}

type MainWnd struct {
	objs   WndObjects
	window *qml.Window
}

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

func run() error {
	// create new qml engine, load qml file
	engine := qml.NewEngine()
	component, err := engine.LoadFile("goConfigAppQt.qml")
	if err != nil {
		return err
	}

	// create MainWnd and qml window
	mainWnd := MainWnd{}
	mainWnd.window = component.CreateWindow(nil)

	// fill mainWnd.objs with qml.Object from the window ('controls', etc)
	err = qmlutil.FetchObjects(&mainWnd.objs, mainWnd.window.Common)
	if err != nil {
		fmt.Println("failed to get objects from window:", err)
		return err
	}

	// hook mainWnd.onQueryButtClicked for query button click
	mainWnd.objs.QueryButton.On("clicked", mainWnd.onQueryButtClicked)
	//engine.Context().SetVar("mainWnd", mainWnd)

	mainWnd.window.Show()
	mainWnd.window.Wait()
	return nil
}

func (w *MainWnd) onQueryButtClicked() {
	w.AddProvider(Provider{"Nexus7", "192.168.1.4"})
}

func (w *MainWnd) AddProvider(p Provider) {
	// create json object string, then call providersMdl.myAppend to append to list
	//##(workaround)
	dicStr := fmt.Sprintf(`{"name": "%v", "address": "%v"}`, p.Name, p.Address)
	w.objs.ProvidersMdl.Call("myAppend", dicStr)
}
