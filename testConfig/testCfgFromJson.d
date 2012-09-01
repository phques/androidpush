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

import json; // probs w. std.json w. GDCm, use a local modif of it!!


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

    void fromJson(FromJson obj) {
        obj.fromJson(this);
    }
}

interface FromJson {
    void fromJson(Wrap json);
}

class LocalRoots : FromJson {
    string videos;
    string music;
    string pictures;
    string download;

    public void fromJson(Wrap json) {
        videos = json.videos.str;
        music = json.music.str;
        pictures = json.pictures.str;
        download = json.download.str;
    }

    override string toString() {
        auto writer = appender!string();
        formattedWrite(writer, "videos=%s\n", videos);
        formattedWrite(writer, "music=%s\n", music);
        formattedWrite(writer, "pictures=%s\n", pictures);
        formattedWrite(writer, "download=%s\n", download);
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
    writeln("\ntestCfg");
    auto json = new Wrap(text, "config");

    Config cfg = new Config;
    json.fromJson(cfg);
    writeln(cfg);

    cfg = new Config;
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
