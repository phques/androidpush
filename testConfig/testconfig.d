module main;

import std.stdio;
import std.file;
import json;

int main(string[] args)
{
    try {
        string text = readText("config.json");
        writeln(text);

        auto json = parseJSON(text);
        writeln(toJSON(&json));

//        auto videos = json["videos"];
        auto obj = json.object;
        writeln(obj);
        writeln(typeid(obj));
        auto videos = obj["videos"];
        writeln(typeid(videos));
        writeln(videos);
    }
    catch (Exception e) {
        writeln("exception : ", e.msg);
    }
	return 0;
}
