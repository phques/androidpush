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

    IupWidget destRootsList;
    IupWidget localRootEdit;
    IupWidget filesList;

    Config config;
    JsonWrap jsonCfg;

    string[] rootsType = ["dummy0", "music", "pictures", "movies", "downloads"];

    string[] droppedFiles;


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
        destRootsList = new IupWidget("destinationRoots");
        localRootEdit = new IupWidget("localRoot");
        filesList = new IupWidget("files");
    }

    static extern(C) int testcb(Ihandle *self, const char* filename, int num, int x, int y) {
        writefln("testcb %s, %s, %s, %s", to!string(filename), num,x,y);
        return IUP_DEFAULT;
    }

    void setupWidgets(){
        setupDestRoots();

//        filesList["DROPTARGET"] = "yes";
        filesList.setDelegate(&this.fileDndCB, "DROPFILES_CB");
//        IupSetCallback(*filesList, "DROPFILES_CB", cast(Icallback)&testcb);
    }

    void setupDestRoots() {
        // fill dest roots list
        foreach (idx, val; rootsType)
            if (idx > 0)
                destRootsList[to!string(idx)] = val;

        // select 1st entry
        destRootsList["VALUE"] = "1";
        setLocalRootFromDestRoot(rootsType[1]);

        destRootsList.setDelegate(&this.rootsCB);
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
        localRootEdit["VALUE"] = newVal;
    }

    // destRoots list callback
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
        filesList["1"] = NULL; // remove all

        foreach (i, filename; droppedFiles) {
            writeln("file: ", filename);
            filesList[to!string(i+1)] = filename;
        }

        // clear droppedFiles
        droppedFiles.length = 0;
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
