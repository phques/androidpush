// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module main;

import std.stdio;
import std.conv;
import std.file;
import std.exception;
import std.variant;
import std.algorithm;

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
    string movies;
    string music;
    string pictures;
    string download;

    void fromJson(Wrap json) {
        movies = json.movies.str;
        music = json.music.str;
        pictures = json.pictures.str;
        download = json.download.str;
    }
}

class Config : FromJson {
    LocalRoots localRoots;

    this() { localRoots = new LocalRoots; }

    void fromJson(Wrap json) {
        localRoots.fromJson(json);
    }
}

class Cfg {
    Variant[string] vals;

    this(string text) {
        auto json = parseJSON(text);
        assert(json.type == JSON_TYPE.OBJECT);
        auto obj = json.object;
        this(obj);
    }

    this(JSONValue[string] obj) {
        foreach(key; obj.keys) {
            auto val = obj[key];
            vals[key] = valToVariant(val);
        }
    }

    static Variant valToVariant(JSONValue val) {
        Variant variant;
        switch(val.type){
            case JSON_TYPE.STRING:
                variant = val.str;
                break;
            case JSON_TYPE.INTEGER:
                variant = val.integer;
                break;
//            case JSON_TYPE.UINTEGER:
//                variant = val.uinteger;
//                break;
            case JSON_TYPE.FLOAT:
                variant = val.floating;
                break;
            case JSON_TYPE.TRUE:
                variant = true;
                break;
            case JSON_TYPE.FALSE:
                variant = false;
                break;
            case JSON_TYPE.NULL:
                variant = null;
                break;
            case JSON_TYPE.ARRAY:
                variant = map!(valToVariant)(val.array);
                break;
            case JSON_TYPE.OBJECT:
                auto tmp = new Cfg(val.object);
                variant = tmp;
                break;
            default:
                throw new Exception("unsupported JSON_TYPE in Cfg.valToVariant");
        }
        return variant;
    }

    Variant opIndex(string k) {
        return vals[k];
    }

//    auto opDispatch(string k)() if (k != null && k.length != 0) {
//        return vals[k];
//    }

}


/*void titi(){
    writeln('\n');

    LocalRoots roots;
    auto ti = typeid(LocalRoots);
    writeln(ti);
    writeln(ti.offTi);

    auto bs = [ __traits(derivedMembers, LocalRoots) ];
    writeln(bs);
}*/

// direct access to JSONValue
void directJson(string text) {
    writeln("\ndirectJson");

    auto json = parseJSON(text);

    auto localRoots = json.object["localRoots"];
    auto movies = localRoots.object["movies"];

    writeln(to!string(movies));
    writeln(movies.str);

    //?? if we dont check the type, we get garbage from the union !
    auto v = movies.integer;
    writeln(v);
}

// using Wrap of JSONValue
void testWrap(string text) {
    writeln("\ntestWrap");

    auto config = new Wrap(text, "config");
    auto localRoots = config.localRoots;
    auto movies = localRoots.movies; //.str;
    writeln(movies.str);
//    auto s = movies.value; // exception
//    writeln(movies.integer); // exception

}

// test Cfg class w. Variants
void testCfg(string text) {
    writeln("\ntestCfg");
    Cfg cfg = new Cfg(text);
    writeln(cfg["localRoots"]);

    auto roots = cfg["localRoots"];
    writeln(to!string(typeid(roots)), roots.type);

    auto r = roots.get!Cfg();
    writeln(r["movies"]);

//    Cfg r2 = roots.get();
//    writeln(r2["movies"]);

//    writeln(cfg.localRoots);
//    writeln(cfg["localRoots"]["movies"]);
//    writeln(cfg.localRoots.movies);
}

void testConfig(string text) {
    writeln("\ntestCfg");
    Config cfg = new Config;
    auto json = new Wrap(text, "config");
    json.fromJson(cfg);
}

int main(string[] args)
{
    try
    {
        string text = readText("config.json");
        writeln(text);

//~         titi();
//        directJson(text);
//        testWrap(text);
//        testCfg(text);
        testConfig(text);
    }
    catch (Exception ex) {
        writeln("ooops exception: ", ex.msg);
    }
	return 0;
}
