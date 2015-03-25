// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module main;

import std.stdio;
import std.conv;
import std.file;
import std.exception;
//import std.variant;
//import std.algorithm;

//import json; // probs w. std.json w. GDCm, use a local modif of it!!
import std.json;
import jsonwrap;


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


int main(string[] args)
{
    try
    {
        string text = readText("config.json");
        writeln(text);
        testWrap(text);
    }
    catch (Exception ex) {
        writeln("ooops exception: ", ex.msg);
    }
	return 0;
}
