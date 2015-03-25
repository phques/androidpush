module main;

// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

import std.stdio;
import std.conv;
import std.socket;
import std.socketstream;
import std.stream;
import std.algorithm;

int main(string[] args)
{
    ushort port = 8888;

    // get server socket port optional param
    if (args.length == 2)
        port = to!ushort(args[1]);

    try
    {
        // connect to server socket
        Socket client = new TcpSocket;
        assert(client.isAlive);
        scope(exit) {
            client.shutdown(SocketShutdown.BOTH);
            client.close();
        }

        writefln("connecting to port %d.", port);
        client.connect(new InternetAddress(port));
        writefln("connected to port %d.", port);

        assert(client.isAlive);

        SocketStream stream = new SocketStream(client, FileMode.In);

        // Get filename, we recv it as 1st line (text)
        // this could cause probs if we dont really receive a line of text !
        // could cause to allocate lots of memory ?
        char[] filename = stream.readLine();
        if (!startsWith(filename, "filename:"))
        {
            writeln("expecting 'filename:..' as 1st line of data!");
            return -1;
        }

        // strip "filename:" to get the filename itself
        filename = filename["filename:".length .. $];
        writeln("filename : ", filename);

        // open outfile file
        BufferedFile outFile = new BufferedFile();
        outFile.open(to!string(filename), FileMode.OutNew);
        scope(exit) outFile.close();

        // receive file data & write to file
        byte[] buffer = new byte[1024*16];
        size_t nbReadBytes;
        size_t total=0;
        do
        {
            // read from socket
            nbReadBytes = stream.readBlock(buffer.ptr, buffer.length);
//            nbReadBytes = client.receive(buffer);
            if (nbReadBytes  > 0)
            {
                // write to file
                total += nbReadBytes;
                writef("%s\r", total);
                outFile.writeBlock(buffer.ptr, nbReadBytes);
            }
        } while (nbReadBytes > 0);

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
