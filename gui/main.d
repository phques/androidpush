// main.d, test using IUP in D (with LED file)

// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

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




class MainWindow : IupWidget {

    this() {
        super("mainDialog");
    }


    // button event callback
    int button2Cb(Ihandle* ihandle) {
        return IUP_DEFAULT;
    }

    void run() {

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

        /* loads LED 'resource' file */
        char* error = IupLoad("androidGUI.led");
        enforce(!error, to!string(error));

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
