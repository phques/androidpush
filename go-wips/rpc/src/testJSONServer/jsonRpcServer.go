package main

import (
	"arith"
	//"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func startServer() {
	arith := new(arith.Arith)

	server := rpc.NewServer()
	server.Register(arith)

	l, e := net.Listen("tcp", ":8222")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	/*for*/ {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func main() {
	startServer()
}
