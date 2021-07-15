package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mox692/gopls_manual_client/client"
	"github.com/mox692/gopls_manual_client/protocol"
	"github.com/sourcegraph/go-langserver/langserver/util"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

var (
	ExitSuccess = 0
	ExitFail    = 1
)

func main() {
	os.Exit(start())
}

func start() int {
	var (
		yamlfile = flag.String("yamlfile", "", "config file")
	)
	flag.Parse()

	// load config
	config, err := client.LoadConfig(*yamlfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitFail
	}
	logger := config.Logger
	logger.Printf("Load config done.\n")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		errCh <- run(config)
	}()

	select {
	case err := <-errCh:
		finish(config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return ExitFail
		}
		logger.Println("Exit Success.")
		return ExitSuccess
	case signal := <-sigCh:
		logger.Printf("Get Signal %d, soon shutdown...\n", signal)
		finish(config)
		return ExitSuccess
	}

}

func run(config *client.Config) error {
	logger := config.Logger
	c, err := net.Dial("tcp", ":"+config.Port)
	if err != nil {
		return err
	}
	handler := client.NewHandler(config)

	conn := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(c, jsonrpc2.VSCodeObjectCodec{}), handler)

	fullpath, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	uri := util.PathToURI(fullpath)
	logger.Printf("Using URI is %s", uri)

	/***************
		initialize
	****************/
	initializesParams := protocol.InitializeParams{
		RootPath:     fullpath,
		RootURI:      protocol.DocumentURI(string(uri)),
		Capabilities: protocol.ClientCapabilities{},
		WorkspaceFolders: []protocol.WorkspaceFolder{
			{
				URI:  "file:///Users/kimuramotoyuki/go/src/github.com/mox692/gopls_manual_client/workspace",
				Name: "dummy_workspace",
			},
		},
	}
	var initializeRes interface{}
	err = conn.Call(context.Background(), "initialize", initializesParams, &initializeRes)
	if err != nil {
		return err
	}

	b, err := json.Marshal(initializeRes)
	if err != nil {
		return err
	}

	logger.Printf("initial Call done, initialize response: %s\n", string(b))

	/***************
		initialized
	****************/
	var initializedParams protocol.InitializedParams
	var initializedRes interface{}
	err = conn.Call(context.Background(), "initialized", initializedParams, &initializedRes)
	if err != nil {
		return err
	}
	b, err = json.Marshal(initializeRes)
	if err != nil {
		return err
	}
	logger.Printf("send initialized Call done, initialized response: %s\n", string(b))

	/***************
		handshake
	****************/
	// next, send handShake method...
	handshakeReq := protocol.HandshakeRequest{}
	handshakeRes := protocol.HandshakeResponse{}
	err = conn.Call(context.Background(), "gopls/handshake", handshakeReq, &handshakeRes)
	if err != nil {
		return err
	}
	logger.Printf("handshake response: %+v\n", handshakeRes)

	/***************
		didOpen
	****************/
	text, err := readFile(config.InitOpenFile)
	if err != nil {
		return err
	}

	didOpenReq := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        protocol.DocumentURI("file:///Users/kimuramotoyuki/go/src/github.com/mox692/gopls_manual_client/workspace/test.go"),
			LanguageID: "go",
			Version:    1,
			Text:       string(text),
		},
	}
	var didOpenRes interface{}
	err = conn.Call(context.Background(), "textDocument/didOpen", didOpenReq, didOpenRes)
	if err != nil {
		return err
	}
	logger.Printf("didOpen success. didOpenres: %+v\n", didOpenRes)

	// TODO: cli interface
	errCh := make(chan error)
	go func() {
		errCh <- startCli(conn, config)
	}()

	// // req didchange to server...
	// var didChangeResult interface{}
	// didChangeParams := protocol.DidChangeTextDocumentParams{
	// 	TextDocument: protocol.VersionedTextDocumentIdentifier{
	// 		Version: 2,
	// 		TextDocumentIdentifier: protocol.TextDocumentIdentifier{
	// 			URI: "dammyURI",
	// 		},
	// 	},
	// }
	// err = conn.Call(context.Background(), "textDocument/didChange", didChangeParams, &didChangeResult)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// logger.Printf("didChange Result: %+v\n", didChangeResult)
	// time.Sleep(time.Second * 100000)
	// fmt.Println("cancel done!! program exit...")

	// wait cli mode
	err = <-errCh
	return err
}

func startCli(conn *jsonrpc2.Conn, config *client.Config) error {
	reqester := client.NewRequester(conn, config)
	var sc = bufio.NewScanner(os.Stdin)
	var err error

	fmt.Println("Start cli mode. Please enter the method name you want to send to the language server.")
	for {
		if sc.Scan() {
			err = reqester.CallMethod(sc.Text(), config)
			if err != nil {
				return err
			}
		}
	}
}

func finish(config *client.Config) {
	// logfileのclose
	config.Logfile.Close()
}

func readFile(pathOrUri string) ([]byte, error) {
	if util.IsURI(lsp.DocumentURI(pathOrUri)) {
		pathOrUri = util.UriToPath(lsp.DocumentURI(pathOrUri))
	}
	bytes, err := ioutil.ReadFile(pathOrUri)
	if err != nil {
		return nil, err
	}
	if len(bytes) > 1024*1024 {
		return nil, errors.New("Read fileSize must be less than 1MB.")
	}
	return bytes, nil
}
