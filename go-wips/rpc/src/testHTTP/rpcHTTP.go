package main

import (
	"arith"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

var listenAddr net.Addr

func startServer() {
	arith := new(arith.Arith)

	server := rpc.NewServer()
	server.Register(arith)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	//http.ListenAndServe(":8222", nil)

	l, e := net.Listen("tcp", ":0")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	listenAddr = l.Addr()
	http.Serve(l, nil)

}

func main() {
	go startServer()
	time.Sleep(time.Duration(100) * time.Millisecond)

	fmt.Println(listenAddr.Network(), listenAddr.String())

	client, err := rpc.DialHTTP(listenAddr.Network(), listenAddr.String())
	//	client, err := rpc.DialHTTP("tcp", ":8222")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	args := &arith.Args{7, 8}
	var reply int

	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

}
