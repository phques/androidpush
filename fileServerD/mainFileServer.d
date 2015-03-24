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
import std.socketstream;
import std.getopt;
//import std.json;
import json;
import std.format;
import std.array;
import std.datetime;

import mimeTypes;

class AndroidPush
{
    private void handleFileParam(string opt, string fileParam)
    {
        // check if it is a file, we also get it's length
        DirEntry entry;
        bool ok = true;
        try {
            entry = dirEntry(fileParam);
            ok = entry.isFile;
        }
        catch (FileException e) {
            ok = false;
        }

        if (ok)
        {
            filePath = fileParam;
            baseFilename = baseName(filePath);
            fileLength = entry.size;
        }
        else
        {
            auto writer = appender!string();
            formattedWrite(writer, "'%s' is not a file", fileParam);

            throw new Exception(writer.data);
        }
    }


    this(string[] args)
    {
        if (args.length == 1)
            throw new Exception("missings arguments");

        getopt(args,
            "androidUdpPort", &androidUdpPort,
            "udpHost", &udpHost,
            "localPort", &localPort,
            "destDirType", &destDirType,
            "destSubdir", &destSubdir,
            "file", &handleFileParam
            );

        if (args.length > 1)
            throw new Exception("invalid argument : " ~ args[1]);
    }

    static void showUsage()
    {
        writeln("usage:");
        writeln("  --androidUdpPort");
        writeln("  --udpHost");
        writeln("  --localPort");
        writeln("  --destDirType");
        writeln("  --destSubdir");
        writeln("  --file");
    }

    void run()
    {
        try
        {
            version(debug_) {
                enum NbRetries = 1;
                enum NbSecWait = 0;
            }
            else {
                enum NbRetries = 15; // 15 x 1sec
                enum NbSecWait = 1;
            }
            bool ok = false;
            for (auto retry=0; retry < NbRetries && !ok; retry++)
            {
                writeln("");
                if (sendUdpPacket())
                {
                    // get 1st connection from android client
                    // quick connect: android client app should connect & send ACK right away
                    waitForAndroidClientSocket("ACK", NbSecWait);
                    if (clientSocket !is null)
                    {
                        scope(exit) {
                            clientSocket.shutdown(SocketShutdown.BOTH);
                            clientSocket.close();
                            clientSocket = null;
                        }

                        // check that we recvd a ACK from android client
                        SocketStream stream = new SocketStream(clientSocket, FileMode.In);
                        char[] ackLine = stream.readLine();
                        char[] pushIdLine = stream.readLine();
                        //stream.close();

                        bool okAck = (ackLine == "ACK");
                        bool okPushId = (pushIdLine == "pushId" ~ to!string(pushId));
                        ok = okAck &&  okPushId;

                        if (ok) {
                            // send the file to client
                            sendFile();
                            writeln("succeeded in sending file ", baseFilename);
                        }
                        else {
                            writeln("did not receive proper ACK & pushId from Android client");
                            debug writefln("ack : '%s'", ackLine);
                            debug writefln("pushId : '%s'", pushIdLine);
                        }
                    }
                }
            }

            if (!ok)
            {
                writeln("\n*** unable to connect to Android client");
            }
        }
        catch (SocketException ex)
        {
            writeln("SocketException ", ex.msg);
        }
        catch (Exception ex)
        {
            writeln("Exception ", ex.msg);
        }
    }

private:

    // create a JSON 'recvFrom' object string
    string makeRecvFromJson(ushort localPort, string destDirType, string subdir,
                            string file, ulong fileLength, int pushId)
    {
        string jsonString =
        `{
            "pushId" : %s,
            "recvFromPort" : %s,
            "destDirType" : "%s",
            "subDir" : "%s",
            "file" : "%s",
            "fileLength" : %s
        }`.format(pushId, localPort, destDirType, subdir, file, fileLength);

        // make sure it parses ! (can throw exception)
        // and get back the actual 'official' JSON string ;-)
        version(MinGW){
            // does not compile Mingw64, gdc.exe (tdm64-1) 4.6.1 !!!
            // missing ref to json skipWhitespace() !!
        }
        else {
            JSONValue val = parseJSON(jsonString);
            //nicer output on android client ;-)
            //jsonString = toJSON(&val);
        }

        return jsonString;
    }

    // send a 'recvFrom' command as a broadcast udp datagram to android client,
    // when it recvs it it then connects to us to get the file
    // ## Throws
    bool sendUdpPacket()
    {
        // open a broadcast udp udpSocket
        Socket udpSocket = new UdpSocket;
        scope(exit) { udpSocket.close(); }

//        auto remoteAddr = new InternetAddress("192.168.1.255", androidUdpPort);
        auto remoteAddr = new InternetAddress(udpHost, androidUdpPort);
        udpSocket.setOption(SocketOptionLevel.SOCKET, SocketOption.BROADCAST, true);

        // build the recvFrom JSON object string to send
        string msg = makeRecvFromJson(localPort, destDirType, destSubdir, baseFilename,
                                      fileLength, pushId);


        // send the recvFrom command to android client
//        writefln("sending %s to %s", msg, to!string(remoteAddr));
        writefln("sending UDP packet to %s", to!string(remoteAddr));
        auto len = udpSocket.sendTo(msg, remoteAddr);

        bool ok = (len > 0 && len != Socket.ERROR);
//        if (ok)
//            writefln("sent datagram (%sbytes)", len);
//        else
        if (!ok)
            writefln("error sending datagram (%s)", len);

        return ok;
    }


    // throws
    void waitForAndroidClientSocket(string waitingFor, int nbSecWaitClientSocket)
    {
        // open listen clientSocket
        Socket listener = new TcpSocket;
        assert(listener.isAlive);
        scope(exit) listener.close();

        auto address = new InternetAddress(localPort);
        listener.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, true);
        listener.bind(address);
        listener.listen(1);
        writefln("Listening on port %d for %s", address.port(), waitingFor);

        // android client should connect right away if it is running & recvs the datagram !
        // so do a select for 1sec waiting for connection
        if (nbSecWaitClientSocket > 0)
        {
            SocketSet set = new SocketSet;
            set.add(listener);

            int selectRes = Socket.select(set, null, null, nbSecWaitClientSocket * 1_000_000);
            if (selectRes < 1) {
                writefln("client failed to connect within %ssec", nbSecWaitClientSocket);
                return;
            }
        }

        // ok we have a connection waiting...
        // accept client connection
        clientSocket = listener.accept();

        writefln("Connection from %s established.", clientSocket.remoteAddress().toString());
        assert(clientSocket.isAlive);
        assert(listener.isAlive);

    }

    void sendFile()
    {
        writefln("file %s", filePath);
        writefln("dest %s/%s", destDirType, destSubdir);

        enum BUFFERSZ = 1024*32;

        StopWatch stopWatch;
        stopWatch.start();

        // open the file to send
        BufferedFile file = new BufferedFile();
        file.open(filePath);
        scope(exit) file.close();

        // sanity check
        if (fileLength != file.size) {
            writefln("fileLength from DirEntry (%s) <> file.size (%s)", fileLength, file.size);
            throw new Exception("fileLength from DirEntry <> file.size");
        }

        // send the file
        byte[] buffer = new byte[BUFFERSZ];
        size_t total=0;

        enum int ProgressBarLen = 50; // nb of characters of progress line
        char[] line = new char[ProgressBarLen];
        double lastPercent=0;

        while (!file.eof())
        {
            // read block from file
            size_t nbBytes = file.readBlock(buffer.ptr, buffer.length);
            if (nbBytes <= 0)
                break;

            // & send the data to the client clientSocket
            clientSocket.send(buffer[0..nbBytes]);

            // display progress : nb bytes total
            total += nbBytes;
            //writef("%s\\%s\r", total, file.size);

            // write out 'progress bar' of percentage
            double percent = total * 100.0 / file.size;
            if (cast(uint)percent != cast(uint)lastPercent) {
                uint nbCharsDone = cast(uint)(percent/100*ProgressBarLen);
                line[0..nbCharsDone] = '=';
                line[nbCharsDone..$] = '-';
                writef(" %02d%% %s\r", cast(int)percent, line);

                stdout.flush();
                lastPercent = percent;
            }
        }

        stopWatch.stop();

        writefln("\n%s\n%s", total, file.size);

        auto duration = stopWatch.peek();
//        writefln("%s:%s", duration.seconds/60, duration.seconds%60);
        writefln("\n%ssecs", duration.seconds);
    }


private:
    ushort androidUdpPort = 4444;
    ushort localPort = 4445;
//    string udpHost = "192.168.1.255"; // broadcast local net .. mine, at home ;-)
    string udpHost = "255.255.255.255"; // broadcast, dont need to know local net
    string destDirType = "";
    string destSubdir = "";
    string filePath;
    string baseFilename;
    ulong fileLength = 0;
    int pushId = 1; // used to identify / synch the 'push' (android client sends it back when it connects)

    Socket clientSocket;

};



int main(string[] args)
{
    AndroidPush androidPush;
    try
    {
        androidPush = new AndroidPush(args);
        androidPush.run();
    }
    catch (Exception ex)
    {
        writeln("ooops, ", ex.msg);
        AndroidPush.showUsage();
    }

    // put breakpt here when in debugger, since terminal wont wait for Enter

	return 0;
}


