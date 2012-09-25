// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module main;

import std.stdio;
import std.conv;
import std.format;
import std.array;
import std.file;
import std.exception;

import std.json;

import kwezd.jsonutil;


//---------------

//-- classes that hold the config
class Config {
    LocalRoots localRoots;
    string lastRoot;
    bool flag;
    float dummy; // if we have an actual entry 'dummy' in json, will assert 'unsupported type'

    this() { localRoots = new LocalRoots; }

    override string toString() {
        auto writer = appender!string();
        formattedWrite(writer,
            "Config{\n%s\n"
            "lastRoot:%s\n"
            "flag:%s\n"
            "}",
            localRoots.toString, lastRoot, flag);
        return writer.data;
    }
}

class LocalRoots {
    string movies;
    string music;
    string pictures;
    string downloads;
    long lnumber;
    int inumber;

    override string toString() {
        auto writer = appender!string();
        formattedWrite(writer,
            "LocalRoots{\n"
            " movies=%s\n"
            " music=%s\n"
            " pictures=%s\n"
            " downloads=%s\n"
            " inumber=%s\n"
            "}",
            movies,music, pictures, downloads, inumber);
        return writer.data;
    }
}

//-------

void testConfig(string text) {

    Config cfg = new Config;
    JsonWrap json = new JsonWrap(text, "config");

    json.Populate(cfg);
    writeln("");
    writeln("movies: ", cfg.localRoots.movies);
    writeln(cfg.toString);
}

int main(string[] args)
{
    try
    {
        string text = readText("config.json");
        writeln(text);
        testConfig(text);
    }
    catch (Exception ex) {
        writeln("ooops exception: ", ex.msg);
    }
	return 0;
}
