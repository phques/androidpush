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


// class that will hold all the controls we work with
// autopopulate with FetchWidgets()
protected class Widgets {
    IupWidget destRootsList;
    IupWidget localRootEdit;
    IupWidget filesList;
};


class MainWindow : IupWidget {

    Widgets w;

    Config config;
    JsonWrap jsonCfg;

    string[] rootsType = ["dummy0", "music", "pictures", "movies", "downloads"];

    string[] droppedFiles;


    this() {
        super("mainDialog");
        w = new Widgets();
        loadConfig();
        FetchWidgets(w);
        setupWidgets();
    }


    version(none)
    static extern(C) int testcb(Ihandle *self, const char* filename, int num, int x, int y) {
        writefln("testcb %s, %s, %s, %s", to!string(filename), num,x,y);
        return IUP_DEFAULT;
    }

    void loadConfig() {
        // create/fill Config config from JSON settings file
        string jsonCfgText = readText("config.json");
        jsonCfg = new JsonWrap(jsonCfgText, "config");

        config = new Config;
        jsonCfg.Populate(config);
    }

    void setupWidgets(){
        setupDestRoots();

//        w.filesList["DROPTARGET"] = "yes";
        w.filesList.setDelegate(&this.fileDndCB, "DROPFILES_CB");
//        IupSetCallback(*w.filesList, "DROPFILES_CB", cast(Icallback)&testcb);
    }

    void setupDestRoots() {
        // fill the destination roots list
        foreach (idx; 1 .. rootsType.length)
            w.destRootsList[to!string(idx)] = rootsType[idx];

        // select 1st entry (need to manually call setLocalRootFromDestRoot)
        w.destRootsList["VALUE"] = "1";
        setLocalRootFromDestRoot(rootsType[1]);

        w.destRootsList.setDelegate(&this.rootsCB);
    }




    // set value of localRootEdit from selected destRoots list item
    // ie: = config[selected='music']
    void setLocalRootFromDestRoot(string itemText) {
        // get the localRoot value for this destinatino root from the config
        auto newVal = jsonCfg.localRoots[itemText].str;
        // set the localRootEdit widget
        w.localRootEdit["VALUE"] = newVal;
    }

    // destRoots list callback
    // calls setLocalRootFromDestRoot when an item is selected
    int rootsCB(Ihandle *ih, char *text, int item, int state) {
        debug writefln("rootscb %s %s %s", to!string(text), item, state);

        if (state == 1) // selected
            setLocalRootFromDestRoot(to!string(text));

        return IUP_DEFAULT;
    }


    // gather all dropped files into droppedFiles[],
    // when last is dropped, process
    int fileDndCB(Ihandle *self, char* text, int num, int x, int y) {
        string filename = to!string(text);

        droppedFiles ~= filename;

        if (num == 0)
            processDroppedFilenames();

        return IUP_DEFAULT;
    }

    void processDroppedFilenames() {
        w.filesList["1"] = NULL; // remove all

        foreach (i, filename; droppedFiles) {
            writeln("file: ", filename);
            w.filesList[to!string(i+1)] = filename;
        }

        // clear droppedFiles
        droppedFiles.length = 0;
    }
}

//---------------------

void loadLed()
{
    /* loads LED 'resource' file */
    version(led) {
        char* error = IupLoad("androidGUI.led");
        enforce(!error, to!string(error));
    }
    else {
        led_load();
    }
}


int main(string[] args)
{

    try {
        /* IUP initialization */
        IupOpenD(args);
        IupControlsOpen() ;
        loadLed();

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
