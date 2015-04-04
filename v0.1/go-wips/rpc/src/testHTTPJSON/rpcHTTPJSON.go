// Note: this uses a 'patched' version net/rpc & net/rpc/json that supports jsonrpc HTTP
package main

import (
	"arith"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

func startServer() {
	arith := new(arith.Arith)

	server := rpc.NewServer()
	server.Register(arith)

	jsonrpc.HandleHTTPDefault(server)
	http.ListenAndServe(":8080", nil)

}

func main() {
	go startServer()
	time.Sleep(time.Duration(100) * time.Millisecond)

	client, err := jsonrpc.DialHTTP("tcp", ":8080")
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

	//##test debug, so we can check http debug
	select {}
}
