// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
module main;

import std.stdio;
import std.file;
import std.path;
import std.string;
import std.conv;
import std.stream;
import std.socket;

int main(string[] args)
{

    string filename;
    ushort port = 8888; // default

    // get filename to send param
    if (args.length < 2)
    {
        writeln("params: filename [port]");
        return -1;
    }
    // also check that it is a valid file
    filename = args[1];
    if (!exists(filename) || !isFile(filename))
    {
        writefln("'%s' is not a file", filename);
        return -1;
    }

    // get optional listen port
    if (args.length == 3)
        port = to!ushort(args[2]);

    try
    {
        // open listen socket
        Socket listener = new TcpSocket;
        assert(listener.isAlive);
        scope(exit) listener.close();

        listener.bind(new InternetAddress(port));
        listener.listen(1);
        writefln("Listening on port %d. file %s", port, filename);

        // accept/waitfor client connection
        Socket sn = listener.accept();
        scope(exit) {
            sn.shutdown(SocketShutdown.BOTH);
            sn.close();
        }

        writefln("Connection from %s established.", sn.remoteAddress().toString());
        assert(sn.isAlive);
        assert(listener.isAlive);

        // send the filename
        sn.send(format("filename:%s\n", baseName(filename)));
version(none) {
        // open the file to send
        BufferedFile file = new BufferedFile();
        file.open(filename);
        scope(exit) file.close();

        // send the file
        byte[] buffer = new byte[1024*16];
        size_t nbReadBytes;
        size_t total=0;
        do
        {
            // read block from file
            nbReadBytes = file.readBlock(buffer.ptr, 1024*16);
            // send the data to the client socket
            sn.send(buffer[0..nbReadBytes]);

            // output nb bytes total
            total += nbReadBytes;
            writef("%s\r", total);

        } while (nbReadBytes > 0);
}
        writeln("");

    }
    catch (SocketException ex)
    {
        writeln("SocketException ", ex.msg);
    }
    catch (Exception ex)
    {
        writeln("Exception ", ex.msg);
    }

//    stdin.readln();

	return 0;
}
