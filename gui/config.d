// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module config;

//-- classes that hold the config
class Config {
    LocalRoots localRoots;

    this() { localRoots = new LocalRoots; }

}

class LocalRoots {
    string movies;
    string music;
    string pictures;
    string downloads;
}

