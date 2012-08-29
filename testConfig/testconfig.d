module main;

import std.stdio;
import std.file;
import json;

int main(string[] args)
{
    try {
        string text = readText("config.json");
        writeln(text);

        // top level object, contains obj localRoots
        auto json = parseJSON(text);
        auto obj = json.object;

        auto localRoots = obj["localRoots"].object;
        auto videos = localRoots["videos"];

        writeln(obj);
        writeln(typeid(obj), '\n');

        writeln(videos, '\n');
        writeln(videos.str);
    }
    catch (Exception e) {
        writeln("exception : ", e.msg);
    }
	return 0;
}
