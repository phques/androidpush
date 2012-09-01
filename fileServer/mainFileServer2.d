// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module main;

import std.stdio;
import std.conv;
import std.socket;
import std.socketstream;
import std.stream;
import std.algorithm;
import std.json;
import std.string;


string makeJson(ushort udpAndroidPort, string dest, string subdir, string file)
{
    string jsonString =
    `{
        "recvFromPort" : %s,
        "destinationType" : "%s",
        "subDir" : "%s",
        "file" : "%s"
    }`.format(udpAndroidPort, dest, subdir, file);

    // make sure it parses ! (can throw exception
    JSONValue val = parseJSON(jsonString);

    // and get back the actual 'official' JSON string ;-)
    jsonString = toJSON(&val);
    debug writeln(jsonString);

    return jsonString;
}



int main(string[] args)
{
    ushort udpAndroidPort = 4445;

    // get server socket udpAndroidPort optional param
    if (args.length == 2)
        udpAndroidPort = to!ushort(args[1]);

    try
    {
        Socket socket = new UdpSocket;
        scope(exit) { socket.close(); }

        // no connect(), broadcast, sendTo()
        auto remoteAddr = new InternetAddress("192.168.1.255", udpAndroidPort);
        socket.setOption(SocketOptionLevel.SOCKET, SocketOption.BROADCAST, true);

        string msg = makeJson(udpAndroidPort, "downloads", "StHilaire-2012-07", "uneimage.jpg");
//        string msg = makeJson(udpAndroidPort, "pictures", "StHilaire-2012-07", "uneimage.jpg");
//        string msg = makeJson(udpAndroidPort, "pictures", "", "uneimage.jpg");

        writefln("sending %s to %s", msg, to!string(remoteAddr));
        auto len = socket.sendTo(msg, remoteAddr);

        if (len == 0 || len == Socket.ERROR)
            writefln("error sending datagram (%s)", len);
        else
            writefln("sent datagram (%s)", len);
    }
    catch (SocketException ex)
    {
        writeln("SocketException ", ex.msg);
    }
    catch (Exception ex)
    {
        writeln("k Exception ", ex.msg);
    }

//    stdin.readln();

	return 0;
}
