module main;

import std.stdio;
import std.conv;
import std.socket;
import std.socketstream;
import std.stream;
import std.algorithm;
import std.json;
import std.string;


string makeJson(ushort port, string dest, string subdir, string file)
{
    string jsonString =
    `{
        "recvFromPort" : %s,
        "destinationType" : "%s",
        "subDir" : "%s",
        "file" : "%s"
    }`.format(port, dest, subdir, file);

    // make sure it parses ! (can throw exception
    JSONValue val = parseJSON(jsonString);

    // and get back the actual 'official' JSON string ;-)
    jsonString = toJSON(&val);
    debug writeln(jsonString);

    return jsonString;
}



int main(string[] args)
{
    ushort port = 4445;

    // get server socket port optional param
    if (args.length == 2)
        port = to!ushort(args[1]);

    try
    {
        Socket socket = new UdpSocket;
        scope(exit) { socket.close(); }

    version(all) {
        // no connect, broadcast, sendTo()
        socket.setOption(SocketOptionLevel.SOCKET, SocketOption.BROADCAST, true);
//        auto remoteAddr = new InternetAddress("192.168.1.123", port);
        auto remoteAddr = new InternetAddress("192.168.1.255", port);

        string msg = makeJson(port, "pictures", "", "uneimage.jpg");

        writefln("sending to %s", to!string(remoteAddr));
        auto len = socket.sendTo(msg, remoteAddr);

        if (len == 0 || len == Socket.ERROR)
            writefln("error sending datagram (%s)", len);
        else
            writefln("sent datagram (%s)", len);
    }
    else {
        // specific address, connect, send()
        auto remoteAddr = new InternetAddress("192.168.1.123", port);
        writefln("connecting to %s", to!string(remoteAddr));
        socket.connect(remoteAddr);

        auto len = socket.send("allo\n");
        if (len == 0 || len == Socket.ERROR)
            writefln("error sending datagram (%s)", len);
        else
            writefln("sent datagram (%s)", len);
    }
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
