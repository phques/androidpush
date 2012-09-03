// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

// main.d, test using IUP in D (with LED file)
module main;

import std.stdio;
import std.string;
import std.exception;
import std.conv;
import std.file;

import iup.iup;
import iup.controls;
import iup.utild;
import iup.widget;
import loadledc;
import kwezd.jsonutil;
import config;


 version = led;

class MainWindow : IupWidget {

    IupWidget destinationRoots;
    IupWidget localRoot;
    Config config;
    JsonWrap jsonCfg;

    string[] rootsType = ["dummy0", "music", "pictures", "movies", "downloads"];

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
        loadConfig();
        getWidgets();
        setupWidgets();
    }

    void getWidgets() {
        // get 'handle' on widgets by name
        destinationRoots = new IupWidget("destinationRoots");
        localRoot = new IupWidget("localRoot");
    }

    void setupWidgets(){
        setupDestRoots();
    }

    void setupDestRoots() {
        // fill dest roots list
        foreach (idx, val; rootsType)
            if (idx > 0)
                destinationRoots[to!string(idx)] = val;

        // select 1st entry
        destinationRoots["VALUE"] = "1";
        setLocalRootFromDestRoot(rootsType[1]);

        destinationRoots.setDelegateSII(&this.rootsCB);
    }

    void loadConfig() {
        string jsonCfgText = readText("config.json");
        jsonCfg = new JsonWrap(jsonCfgText, "config");

        config = new Config;
        jsonCfg.Populate(config);
    }



    // set value of localRoot edit from selected destRoots list item
    // ie: = config[selected='music']
    void setLocalRootFromDestRoot(string itemText) {
        auto newVal = jsonCfg.localRoots[itemText].str;
        localRoot["VALUE"] = newVal;
    }

    // destRoots list callback
    int rootsCB(Ihandle *ih, char *text, int item, int state) {
        debug writefln("rootscb %s %s %s", to!string(text), item, state);

        if (state == 1) // selected
            setLocalRootFromDestRoot(to!string(text));

        return IUP_DEFAULT;
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

        /* shows dialog */
        window.Show();

        /* main loop */
        IupMainLoop();
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
