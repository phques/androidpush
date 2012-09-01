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

//import json; // probs w. std.json w. GDCm, use a local modif of it!!
import std.json;



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

// test Cfg class w. Variants
void testCfg(string text) {
    writeln("\ntestCfg");
    Cfg cfg = new Cfg(text);
    writeln(cfg["localRoots"]);

    auto roots = cfg["localRoots"];
    writeln(to!string(typeid(roots)), roots.type);

    auto r = roots.get!Cfg();
    writeln(r["videos"]);

//    Cfg r2 = roots.get();
//    writeln(r2["videos"]);

//    writeln(cfg.localRoots);
//    writeln(cfg["localRoots"]["videos"]);
//    writeln(cfg.localRoots.videos);
}


int main(string[] args)
{
    try
    {
        string text = readText("config.json");
        writeln(text);
        testCfg(text);
    }
    catch (Exception ex) {
        writeln("ooops exception: ", ex.msg);
    }
	return 0;
}
