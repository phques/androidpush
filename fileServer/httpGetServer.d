
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
import std.json;
import std.format;
import std.array;

import mimeTypes;

class AndroidHttpPush
{
    private void handleFileParam(string opt, string fileParam)
    {
        filePath = fileParam;

        if (!exists(filePath) || !isFile(filePath))
        {
            auto writer = appender!string();
            formattedWrite(writer, "'%s' is not a file", filePath);

            throw new Exception(writer.data);
        }

        baseFilename = baseName(filePath);
    }


    this(string[] args)
    {
        if (args.length == 1)
            throw new Exception("missings arguments");

        getopt(args,
            "androidUdpPort", &androidUdpPort,
            "localHttpPort", &localHttpPort,
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
        writeln("  --localHttpPort");
        writeln("  --destDirType");
        writeln("  --destSubdir");
        writeln("  --file");
    }

    void run()
    {
        try
        {
            enum NbRetries = 15; // 15 x 1sec
            bool ok = false;
            for (auto retry=0; retry < NbRetries && !ok; retry++)
            {
                writeln("");
                if (sendUdpPacket())
                {
                    // get 1st connection from android client
                    // quick connect: android client app should connect & send ACK right away
                    waitForAndroidClientSocket("ACK", 1);
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
                        stream.close();

                        ok = ackLine == "ACK";
                    }
                }
            }

            if (ok)
            {
                // get connection to android client
                writeln("");
                prepareClientSocket();
                if (clientSocket !is null)
                {
                    scope(exit) {
                        clientSocket.shutdown(SocketShutdown.BOTH);
                        clientSocket.close();
                        clientSocket = null;
                    }

                    // send the file to client
                    sendFile();
                    writeln("succeeded in sending file ", baseFilename);
                }
            }
            else
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
    ushort androidUdpPort = 4445;
    ushort localHttpPort = 8080; // 80 not allowed on my ubuntu ;-p
    string destDirType = "downloads";
    string destSubdir = "";
    string filePath;
    string baseFilename;

    Socket clientSocket;

    // create a JSON 'recFrom' object string
    string makeRecvFromJson(ushort localPort, string destDirType, string subdir, string file)
    {
        string jsonString =
        `{
            "recvFromPort" : %s,
            "destDirType" : "%s",
            "subDir" : "%s",
            "file" : "%s"
        }`.format(localPort, destDirType, subdir, file);

        // make sure it parses ! (can throw exception)
        // and get back the actual 'official' JSON string ;-)
        JSONValue val = parseJSON(jsonString);
        jsonString = toJSON(&val);

        return jsonString;
    }

    // send a 'recvFrom' command as a broadcast udp datagram to android client,
    // when it recvs it it then connects to us to get the file (HTTP GET)
    // ## Throws
    bool sendUdpPacket()
    {
        // open a broadcast udp udpSocket
        Socket udpSocket = new UdpSocket;
        scope(exit) { udpSocket.close(); }

        auto remoteAddr = new InternetAddress("192.168.1.255", androidUdpPort);
        udpSocket.setOption(SocketOptionLevel.SOCKET, SocketOption.BROADCAST, true);

        // build the recvFrom JSON object string to send
        string msg = makeRecvFromJson(localHttpPort, destDirType, destSubdir, baseFilename);

        // send the recvFrom command to android client
        //writefln("sending %s to %s", msg, to!string(remoteAddr));
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

        auto address = new InternetAddress(localHttpPort);
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


    void readClientData()
    {
            SocketStream stream = new SocketStream(clientSocket, FileMode.In);

            foreach(ulong n, char[] line; stream)
            {
                writefln("=> %s", line);
                if (line == "") // end of headers
                    break;
            }

    }

    // throws
    void prepareClientSocket()
    {
        // get connection from android client,
        // here we wait until android download mgr service connects
        waitForAndroidClientSocket("HTTP GET", 0);

        if (clientSocket !is null)
        {
            // seems to need this !!

            // Read the incoming GET request (need to read, otherwise it doesnt work !)
            // ##TODO we should interpret the GET http command !!
            // ## check that it is really the android client, asking for the correct file
            readClientData();
        }
    }


    void sendFile()
    {
        writefln("file %s\n", filePath);

        enum BUFFERSZ = 1024*32;

        // open the file to send
        BufferedFile file = new BufferedFile();
        file.open(filePath);
        scope(exit) file.close();

        string fileSizeStr = to!string(file.size);

        // send back headers
        debug writeln("Content-Length: " ~ fileSizeStr ~ "\n");

        clientSocket.send("HTTP/1.0 200 OK\n");
        clientSocket.send("Server: androidHttpPush\n");
        clientSocket.send("Connection: close\n");
        clientSocket.send("Content-type: " ~ MimeTypes.getFromFile(baseFilename) ~ "\n");
//        clientSocket.send("Content-type: application/octet-stream\n");
//        clientSocket.send("Content-type: image/jpeg\n");
        clientSocket.send("Content-Length: " ~ fileSizeStr ~ "\n");
        clientSocket.send("\n");


        // send the file
        byte[] buffer = new byte[BUFFERSZ];
        size_t total=0;
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
            writef("%s\\%s\r", total, fileSizeStr);
        }

        writeln("");

        readClientData(); // seems to need this !!
    }


};



int main(string[] args)
{
    AndroidHttpPush androidPush;
    try
    {
        androidPush = new AndroidHttpPush(args);
        androidPush.run();
    }
    catch (Exception ex)
    {
        writeln("ooops, ", ex.msg);
        AndroidHttpPush.showUsage();
    }

    // when in debugger
    stdin.readln();
	return 0;
}


