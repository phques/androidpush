// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module jsonwrap;

import std.stdio;
import std.exception;
import std.json;


class Wrap {
    JSONValue json;
    string name;

    this(string jsonString, string name) { this.name = name; this.json = parseJSON(jsonString); }
    this(JSONValue json, string name) { this.name = name; this.json = json; }

    // shortcut objWrap.memberX = objWrap.json.object["memberX"]
    Wrap opDispatch(string part)() {
        return opIndex(part);
    }

    // objWrap["memberX"] = objWrap.memberX = objWrap.json.object["memberX"]
    Wrap opIndex(string part) {
        string fullname = name ~ '.' ~ part;
        enforce(json.type == JSON_TYPE.OBJECT, "Expecting JSON.OBJECT, to access " ~ fullname);
        enforce(part in json.object, "Missing json member: " ~ fullname);

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
    bool boolean() {
        switch (json.type) {
            case JSON_TYPE.TRUE: return true;
            case JSON_TYPE.FALSE: return false;
            default: break;
        }
        enforce(false, "Expecting JSON.TRUE/FALSE, for " ~ name);
        return false; // keep compiler from complaining ;-)
    }

    JSONValue[string] object() {
        enforce(json.type == JSON_TYPE.OBJECT, "Expecting JSON.OBJECT, for " ~ name);
        return json.object;
    }
}

