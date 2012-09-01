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

//import json; // probs w. std.json w. GDCm, use a local modif of it!!
import std.json;
import jsonwrap;


interface FromJson {
    void fromJson(Wrap json);
}

class LocalRoots : FromJson {
    string movies;
    string music;
    string pictures;
    string downloads;

    public void fromJson(Wrap json) {
        movies = json.movies.str;
        music = json.music.str;
        pictures = json.pictures.str;
        downloads = json.downloads.str;
    }

    override string toString() {
        auto writer = appender!string();
        formattedWrite(writer, "movies=%s\n", movies);
        formattedWrite(writer, "music=%s\n", music);
        formattedWrite(writer, "pictures=%s\n", pictures);
        formattedWrite(writer, "downloads=%s\n", downloads);
        return writer.data;
    }
}

class Config : FromJson {
    LocalRoots localRoots;

    this() { localRoots = new LocalRoots; }

    void fromJson(Wrap json) {
        localRoots.fromJson(json.localRoots);
    }

    override string toString() {
        auto writer = appender!string();
        formattedWrite(writer, "Config{\n%s\n}", localRoots.toString);
        return writer.data;
    }
}


void testConfig(string text) {
    auto json = new Wrap(text, "config");

    Config cfg = new Config;
    cfg.fromJson(json);
    writeln(cfg);

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
