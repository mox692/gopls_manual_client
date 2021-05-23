package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

type clientHandler struct {
}

type InitializeParams struct {
	RootPath string
	RootURI  string
}

func (c *clientHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	fmt.Printf("called!! request is : %+v\n", req)
}

func main() {
	c, err := net.Dial("tcp", ":1234")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	conn := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(c, jsonrpc2.VSCodeObjectCodec{}), &clientHandler{})

	fmt.Printf("%+v\n", conn)
	done := make(chan bool)

	// fullpath, err := filepath.Abs(".")
	// if err != nil {
	// 	panic(err)
	// }

	params := InitializeParams{
		RootPath: "sfsfs",
	}

	var got interface{}
	err = conn.Call(context.Background(), "initialize", params, &got)
	if err != nil {
		panic(err)
	}

	fmt.Println("Call done!!!")
	go func() {
		<-conn.DisconnectNotify()
		done <- true
	}()

}
