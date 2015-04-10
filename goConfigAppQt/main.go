package main

import (
	"fmt"
	"os"

	"gopkg.in/qml.v1"
)

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

var providers *Providers

func run() error {
	engine := qml.NewEngine()
	component, err := engine.LoadFile("goConfigAppQt.qml")
	if err != nil {
		return err
	}

	window := component.CreateWindow(nil)

	providers = &Providers{}
	providers.obj = window.ObjectByName("providersMdl")
	//	engine.Context().SetVar("providers", providers)

	queryButton := window.ObjectByName("queryButton")
	queryButton.On("clicked", onQueryButtClicked)

	window.Show()
	window.Wait()
	return nil
}

func onQueryButtClicked() {
	providers.Add(Provider{"Nexus7", "192.168.1.4"})
}

type Provider struct {
	Name    string
	Address string
}

type Providers struct {
	obj qml.Object
}

func (provs *Providers) Add(p Provider) {
	// create json object string, then call myAppend to append to list
	dicStr := fmt.Sprintf(`{"name": "%v", "address": "%v"}`, p.Name, p.Address)
	provs.obj.Call("myAppend", dicStr)
}
