package main

import (
	"fmt"
	"github.com/grd/iup"
	"iuputil"
	"time"
)

type MyWidgets struct {
	MainDialog *iup.Ihandle `IUP:"mainDialog"`
	LocalRoot  *iup.Ihandle `IUP:"localRoot"`
	Files      *iup.Ihandle
}

var idleChan chan string
var cpt int
var myWidgets MyWidgets

func idleFunc() int {
	select {
	case cmd := <-idleChan:
		iup.SetAttribute(myWidgets.LocalRoot, "VALUE", cmd)
		fmt.Println("got something to do in idle: ", cmd)
	case <-time.After(time.Duration(150 * time.Millisecond)):
	}
	return iup.DEFAULT
}

func main() {
	iup.Open()
	defer iup.Close()

	if err := iup.Load("androidGUI.led"); err != nil {
		fmt.Println(err)
		return
	}

	if err := iuputil.FetchWidgets(&myWidgets); err != nil {
		fmt.Println("FetchWidgets failed : ", err)
		return
	}

	/* shows dialog */
	iup.Show(myWidgets.MainDialog)

	// prepare a channel for the idle callback,
	// start a goroutine to send a msg on the channel after some time
	idleChan = make(chan string)

	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		idleChan <- "command to process by UI thread"
	}()

	// hook our idle func
	iup.SetIdleFunc(idleFunc)

	/* main loop */
	iup.MainLoop()

}
