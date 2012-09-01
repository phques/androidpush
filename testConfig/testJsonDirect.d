// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module main;

import std.stdio;
import std.file;
import std.exception;
import std.conv;

//import json; // probs w. std.json w. GDCm, use a local modif of it!!
import std.json;



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


int main(string[] args)
{
    try
    {
        string text = readText("config.json");
        writeln(text);
        directJson(text);
    }
    catch (Exception ex) {
        writeln("ooops exception: ", ex.msg);
    }
	return 0;
}
