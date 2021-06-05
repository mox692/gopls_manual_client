package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/mox692/gopls_manual_client/client"
	"github.com/mox692/gopls_manual_client/protocol"
	"github.com/sourcegraph/go-langserver/langserver/util"
	"github.com/sourcegraph/jsonrpc2"
)

/**
 * these struct is referenced from  golang/tools/internal/lsp/protocol/tsprotocol.go
 * ref: https://github.com/golang/tools/blob/master/internal/lsp/protocol/tsprotocol.go
 */
// type clientHandler struct {
// }

// func (c *clientHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
// 	fmt.Printf("called!! request is : %+v\n", req)
// 	if !req.Notif {
// 		return
// 	}

// 	switch req.Method {
// 	case "textDocument/documentHighlight":
// 		var params lsp.PublishDiagnosticsParams
// 		b, _ := req.Params.MarshalJSON()
// 		json.Unmarshal(b, &params)
// 	}
// }

func main() {
	var logger *log.Logger
	var (
		port    = flag.String("port", "37375", "gopls's port")
		logfile = flag.String("logfile", "", "logfile")
	)
	flag.Parse()

	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		logger = log.New(f, "", log.LstdFlags)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	c, err := net.Dial("tcp", ":"+*port)
	if err != nil {
		logger.Fatal(err)
	}

	config := client.LoadConfig()
	handler := client.NewHandler(config)

	conn := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(c, jsonrpc2.VSCodeObjectCodec{}), handler)

	fullpath, err := filepath.Abs(".")
	if err != nil {
		logger.Fatal(err)
	}

	uri := util.PathToURI(fullpath)

	params := protocol.InitializeParams{
		RootPath: fullpath,
		RootURI:  protocol.DocumentURI(string(uri)),
	}

	var initializeRes interface{}
	err = conn.Call(context.Background(), "initialize", params, &initializeRes)
	if err != nil {
		logger.Fatal(err)
	}

	b, err := json.Marshal(initializeRes)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("initial Call done, initialize response: %s\n", string(b))

	// next, send handShake method...
	handshakeReq := protocol.HandshakeRequest{}
	handshakeRes := protocol.HandshakeResponse{}
	err = conn.Call(context.Background(), "gopls/handshake", handshakeReq, &handshakeRes)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("handshake response: %+v\n", handshakeRes)

	// wait for message from server...
	ctx, cancel := context.WithCancel(context.Background())
	go waitReq(ctx)

	// req didOpen to server...
	didOpenReq := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "",
			LanguageID: "go",
			Version:    1,
			Text: `
			package main

			func main() {
				fmt.Println("hello 世界")
			}
			`,
		},
	}
	var didOpenRes interface{}
	err = conn.Call(context.Background(), "textDocument/didOpen", didOpenReq, didOpenRes)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("didOpen success. didOpenres: %+v\n", didOpenRes)

	// req didchange to server...
	var didChangeResult interface{}
	didChangeParams := protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			Version: 2,
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{
				URI: "dammyURI",
			},
		},
	}
	err = conn.Call(context.Background(), "textDocument/didChange", didChangeParams, &didChangeResult)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("didChange Result: %+v\n", didChangeResult)
	cancel()
	time.Sleep(time.Second * 100000)
	fmt.Println("cancel done!! program exit...")
}

func waitReq(ctx context.Context) {
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("fdsafs")
		select {
		case <-ctx.Done():
			fmt.Printf("waitReq Done...\n")
			return
		default:
		}
	}
}
