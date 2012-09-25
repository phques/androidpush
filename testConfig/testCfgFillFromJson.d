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

//-----------

version = wrap;

// Functions to fill an object from a 'parallel' json object
//


version(wrap)

// by using class Wrap, we get 'auto type safety' & 'qualified name' for errors
// nb: this could be nice as a method of Wrap class
void Fill(Class)(ref Class obj, Wrap json) {

    json.object(); // enforce(type=object) !

    // for each member of the class
    foreach (m; __traits(derivedMembers, Class)) {

        // do we have a member with same name in json objet ?
        // if so, assign according to type of class member
        if (m in json.object) {
            static if (is(typeof(__traits(getMember, obj, m)) == string)) {
                __traits(getMember, obj, m) = json[m].str;
            }
            else static if (is(typeof(__traits(getMember, obj, m)) == int)) {
                __traits(getMember, obj, m) = cast(int)json[m].integer;
            }
            else static if (is(typeof(__traits(getMember, obj, m)) == long)) {
                __traits(getMember, obj, m) = json[m].integer;
            }
            else static if (is(typeof(__traits(getMember, obj, m)) == bool)) {
                __traits(getMember, obj, m) = json[m].boolean;
            }
            else static if (is(typeof(__traits(getMember, obj, m)) == class)) {
                Fill(__traits(getMember, obj, m), json[m]);  // object, recurse
            }
            else {
                assert(false, "Unsupported type " ~ typeof(__traits(getMember, obj, m)).stringof);
            }
        }
    }
}

else version(nowrap)

// direct use of JSONValue
// this version would need manual checks on type of json value
// (w. ugly debug writeln()s ;-p)
void Fill(Class)(ref Class obj, JSONValue json) {

    debug writeln("fill: ", [__traits(derivedMembers, Class)]);
    debug writeln("with: ", to!string(json.object));

    // for each member of the class
    foreach (m; __traits(derivedMembers, Class)) {

        debug writeln(typeof(__traits(getMember, obj, m)).stringof, " ", m);

        // do we have a member with same name in json objet ?
        // if so, assign according to type of class member
        //## need check of jsonvalue type !!
        if (m in json.object) {
            static if (is(typeof(__traits(getMember, obj, m)) == string)) {
                debug writeln("  assigning str to ", m);
                __traits(getMember, obj, m) = json.object[m].str;
            }
            else static if (is(typeof(__traits(getMember, obj, m)) == int)) {
                debug writeln("  assigning int to ", m);
                __traits(getMember, obj, m) = cast(int)json.object[m].integer;
            }
            else static if (is(typeof(__traits(getMember, obj, m)) == long)) {
                debug writeln("  assigning int to ", m);
                __traits(getMember, obj, m) = json.object[m].integer;
            }
            else static if (is(typeof(__traits(getMember, obj, m)) == bool)) {
                debug writeln("  assigning bool to ", m);
                __traits(getMember, obj, m) = json.object[m].type == JSON_TYPE.TRUE;
            }
            else static if (is(typeof(__traits(getMember, obj, m)) == class)) {
                debug writeln("  recursing on object");
                Fill(__traits(getMember, obj, m), json.object[m]);
            }
            else {
                debug writeln("  unsupported type ");
            }
        }
        else {
            // nb: outputs for every method/CTOR etc...!
            debug writeln("  warning: no corresponding JSON member '", m, "'");
        }
    }
}

//---------------

//-- classes that hold the config
class Config {
    LocalRoots localRoots;
    string lastRoot;
    float dummy; // if we have an actual entry 'dummy' in json, will assert 'unsupported type'
    bool flag;

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
    version(nowrap) auto json = parseJSON(text);
    else            auto json = new Wrap(text, "config");

    Fill(cfg, json);
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
