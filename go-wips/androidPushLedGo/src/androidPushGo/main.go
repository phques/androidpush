package main

import (
	"fmt"
	"github.com/grd/iup"
	"iuputil"
	"runtime"
	"time"
)

type MyWidgets struct {
	MainDialog *iup.Ihandle `IUP:"mainDialog"`
	LocalRoot  *iup.Ihandle `IUP:"localRoot"`
	Files      *iup.Ihandle
	Push       *iup.Ihandle `IUP:"pushButton"`
}

var cmdChan chan string
var quitChan chan struct{}
var cpt int
var myWidgets MyWidgets

//---------

func idleFunc1() int {
	select {
	case cmd := <-cmdChan:
		iup.SetAttribute(myWidgets.LocalRoot, "VALUE", cmd)
		fmt.Println("got something to do in idle: ", cmd)

	case <-time.After(time.Duration(50 * time.Millisecond)):
	}
	return iup.DEFAULT
}

//----------

func idleFunc2() int {
	var i = 0
	for iup.LoopStepWait() != iup.CLOSE {
		fmt.Println("iup.LoopStepWait", i)
		i++

		// this won't get called until an UI event occurs !
		// on Linux/GTK we do get called (some timer event ?)
		select {
		case cmd := <-cmdChan:
			iup.SetAttribute(myWidgets.LocalRoot, "VALUE", cmd)
			fmt.Println("got something to do in idle: ", cmd)
		default:
		}
	}

	fmt.Println("idleFunc2 exit")
	return iup.CLOSE
}

//-----

// ## does not work !!
// the GUI is frozen .. seems like iup.LoopStepWait() never returns

func myIupLoop() {
	runtime.LockOSThread()

	var i = 0
	for iup.LoopStepWait() != iup.CLOSE {
		fmt.Println("iup.LoopStepWait", i)
		i++
	}

	close(quitChan)
}

func idleFunc3() int {

	//runtime.LockOSThread()

	go myIupLoop()

	for {
		select {
		case cmd := <-cmdChan:
			iup.SetAttribute(myWidgets.LocalRoot, "VALUE", cmd)
			fmt.Println("got something to do in idle: ", cmd)

		case <-quitChan:
			break
		}
	}

	fmt.Println("idleFunc3 exit")
	return iup.CLOSE
}

//----------

func main() {
	runtime.LockOSThread()

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
	cmdChan = make(chan string)
	quitChan = make(chan struct{})

	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		cmdChan <- "command to process by UI thread"
	}()

	// hook our idle func
	iup.SetIdleFunc(idleFunc1)

	/* main loop */
	iup.MainLoop()

}
