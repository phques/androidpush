module main;

import std.stdio;
import std.conv;
import std.file;
import std.traits;

import json; // probs w. std.json !!

class Wrap {
    JSONValue json;

    this(string jsonString) { this.json = parseJSON(jsonString); }
    this(JSONValue json) { this.json = json; }

    auto opDispatch(string name)() {
        debug writeln("opdispatch ", name);
        assert(json.type == JSON_TYPE.OBJECT);
        auto v = mixin(`json.object["`~name~`"]`);
        return new Wrap(v);
    }

    string str() {
        assert(json.type == JSON_TYPE.STRING);
        return json.str;
    }
    long integer() {
        assert(json.type == JSON_TYPE.INTEGER);
        return json.integer;
    }

}

struct LocalRoots {
    string videos;
    string music;
    string pictures;
    string download;
}

void toto(string text) {
    writeln('\n');

    auto json = parseJSON(text);
    auto obj = json.object;

    auto localRoots = obj["localRoots"].object;
    auto videos = localRoots["videos"];

    writeln(to!string(obj));
    writeln(typeid(obj), '\n');

    writeln(to!string(videos));
    writeln(videos.str);
}

void tata(string text) {
    writeln('\n');

    auto config = new Wrap(text);
    auto localRoots = config.localRoots;
    auto videos = localRoots.videos.str;
    writeln(videos);
//    auto s = videos.value; // asserts
//    writeln(videos.integer); // asserts

}

int main(string[] args)
{
    LocalRoots roots;
    auto ti = typeid(LocalRoots);
    writeln(ti);
    writeln(ti.offTi);

    auto bs = [ __traits(derivedMembers, LocalRoots) ];
    writeln(bs);

    string text = readText("config.json");
    writeln(text);

    toto(text);
    tata(text);

	return 0;
}
