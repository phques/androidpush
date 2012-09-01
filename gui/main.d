// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// main.d, test using IUP in D (with LED file)
module main;

import std.stdio;
import std.string;
import std.utf;
import std.exception;
import std.conv;
import std.typecons;

import iup.iup;
import iup.controls;
import iup.utild;
import iup.widget;
import loadledc;


 version = led;

class MainWindow : IupWidget {

    this() {
        /* loads LED 'resource' file */
        version(led) {
            char* error = IupLoad("androidGUI.led");
            enforce(!error, to!string(error));
        }
        else {
            led_load();
        }

        super("mainDialog");
    }


    // button event callback
    int button2Cb(Ihandle* ihandle) {
        return IUP_DEFAULT;
    }

    void run() {

        //## test
        auto list = new IupWidget("destinationRoot");
        list["ACTIVE"] = "no";
        writeln(*list);

        /* shows dialog */
        this.Show();

        /* main loop */
        IupMainLoop();
    }
}

//---------------------

int main(string[] args)
{

    try {
        /* IUP initialization */
        IupOpenD(args);
        IupControlsOpen() ;

        scope auto window = new MainWindow;
        window.run();

    }
    catch (Exception e) {
        IupMessage("error", e.msg.toStringz);
    }
    finally {
        /* ends IUP */
        IupControlsClose();
        IupClose();
    }

	return 0;
}
