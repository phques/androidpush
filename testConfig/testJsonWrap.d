// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module main;

import std.stdio;
import std.conv;
import std.file;
import std.exception;
//import std.variant;
//import std.algorithm;

//import json; // probs w. std.json w. GDCm, use a local modif of it!!
import std.json;


class Wrap {
    JSONValue json;
    string name;

    this(string jsonString, string name) { this.name = name; this.json = parseJSON(jsonString); }
    this(JSONValue json, string name) { this.name = name; this.json = json; }

    // shortcut objWrap.memberX = objWrap.json.object["memberX"]
    auto opDispatch(string part)() {
        string fullname = name ~ '.' ~ part;

        debug writeln("opdispatch ", fullname);
        enforce(json.type == JSON_TYPE.OBJECT, "Expecting JSON.OBJECT, for " ~ fullname);

        auto v = json.object[part];
        return new Wrap(v, fullname);
    }

    string str() {
        enforce(json.type == JSON_TYPE.STRING, "Expecting JSON.STRING, for " ~ name);
        return json.str;
    }
    long integer() {
        enforce(json.type == JSON_TYPE.INTEGER, "Expecting JSON.INTEGER, for " ~ name);
        return json.integer;
    }

}


// using Wrap of JSONValue
void testWrap(string text) {
    writeln("\ntestWrap");

    auto config = new Wrap(text, "config");
    auto localRoots = config.localRoots;
    auto videos = localRoots.videos; //.str;
    writeln(videos.str);
//    auto s = videos.value; // exception
//    writeln(videos.integer); // exception

}


int main(string[] args)
{
    try
    {
        string text = readText("config.json");
        writeln(text);
        testWrap(text);
    }
    catch (Exception ex) {
        writeln("ooops exception: ", ex.msg);
    }
	return 0;
}
