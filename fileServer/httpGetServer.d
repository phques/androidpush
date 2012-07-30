
module main;

import std.stdio;
import std.file;
import std.path;
import std.string;
import std.conv;
import std.stream;
import std.socket;
import std.socketstream;

int main(string[] args)
{

    string filename;
    ushort port = 80; // default

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

        listener.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, true);

        listener.bind(new InternetAddress(port));
        listener.listen(1);
        writefln("Listening on port %d. file %s", port, filename);

        // accept/waitfor client connection
        Socket socket = listener.accept();
//        scope(exit) {
//            socket.shutdown(SocketShutdown.BOTH);
//            socket.close();
//        }

        writefln("Connection from %s established.", socket.remoteAddress().toString());
        assert(socket.isAlive);
        assert(listener.isAlive);

        // ##TODO we should interpret the GET http command !!

        SocketStream stream = new SocketStream(socket, FileMode.In);

         foreach(ulong n, char[] line; stream)
         {
             writefln("line : %s", line);
             if (line == "") // end of headers
                break;
         }


        // open the file to send
        BufferedFile file = new BufferedFile();
        file.open(filename);
        scope(exit) file.close();

        // send back headers
        debug writeln("Content-Length: " ~ to!string(file.size) ~ "\n");

        socket.send("HTTP/1.0 200 OK\n");
        socket.send("Server: kwezMiniHttp\n");
        socket.send("Connection: close\n");
//        socket.send("Content-type: application/octet-stream\n");
        socket.send("Content-type: image/jpeg\n");
        socket.send("Content-Length: " ~ to!string(file.size) ~ "\n");
        socket.send("\n");


version(all) {
        // send the file
        byte[] buffer = new byte[1024*16];
        size_t nbReadBytes;
        size_t total=0;
        do
        {
            // read block from file
            nbReadBytes = file.readBlock(buffer.ptr, 1024*16);
            // send the data to the client socket
            socket.send(buffer[0..nbReadBytes]);

            // output nb bytes total
            total += nbReadBytes;
            writef("%s\r", total);

        } while (nbReadBytes > 0);

        writeln("");
}
        foreach(ulong n, char[] line; stream)
         {
             writefln("line : %s", line);
//             if (line == "") // end of headers
//                break;
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

//    stdin.readln();

	return 0;
}
