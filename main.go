package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/mox692/gopls_manual_client/protocol"
	"github.com/sourcegraph/go-langserver/langserver/util"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

/**
 * these struct is referenced from  golang/tools/internal/lsp/protocol/tsprotocol.go
 * ref: https://github.com/golang/tools/blob/master/internal/lsp/protocol/tsprotocol.go
 */
type clientHandler struct {
}

type InitializeParams struct {
	RootPath string
	RootURI  lsp.DocumentURI
}

type DocumentHighlightParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

type PartialResultParams struct {
	/**
	 * An optional token that a server can use to report partial results (e.g. streaming) to
	 * the client.
	 */
	PartialResultToken interface{} `json:"partialResultToken,omitempty"`
}

type WorkDoneProgressParams struct {
	/**
	 * An optional token that a server can use to report work done progress.
	 */
	WorkDoneToken interface{} `json:"workDoneToken,omitempty"`
}

type TextDocumentPositionParams struct {
	/**
	 * The text document.
	 */
	TextDocument TextDocumentIdentifier `json:"textDocument,string"` // uri
	/**
	 * The position inside the text document.
	 */
	Position Position `json:"position"`
}

type Position struct {
	/**
	 * Line position in a document (zero-based).
	 */
	Line uint32 `json:"line"`
	/**
	 * Character offset on a line in a document (zero-based). Assuming that the line is
	 * represented as a string, the `character` value represents the gap between the
	 * `character` and `character + 1`.
	 *
	 * If the character value is greater than the line length it defaults back to the
	 * line length.
	 */
	Character uint32 `json:"character"`
}

/**
 * A literal to identify a text document in the client.
 */
type TextDocumentIdentifier struct {
	/**
	 * The text document's uri.
	 */
	URI string `json:"uri,string"`
}

func (c *clientHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	fmt.Printf("called!! request is : %+v\n", req)
	if !req.Notif {
		return
	}

	switch req.Method {
	case "textDocument/documentHighlight":
		var params lsp.PublishDiagnosticsParams
		b, _ := req.Params.MarshalJSON()
		json.Unmarshal(b, &params)
	}
}

type handshakeRequest struct {
	// ServerID is the ID of the server on the client. This should usually be 0.
	ServerID string `json:"serverID"`
	// Logfile is the location of the clients log file.
	Logfile string `json:"logfile"`
	// DebugAddr is the client debug address.
	DebugAddr string `json:"debugAddr"`
	// GoplsPath is the path to the Gopls binary running the current client
	// process.
	GoplsPath string `json:"goplsPath"`
}
type handshakeResponse struct {
	// SessionID is the server session associated with the client.
	SessionID string `json:"sessionID"`
	// Logfile is the location of the server logs.
	Logfile string `json:"logfile"`
	// DebugAddr is the server debug address.
	DebugAddr string `json:"debugAddr"`
	// GoplsPath is the path to the Gopls binary running the current server
	// process.
	GoplsPath string `json:"goplsPath"`
}

func main() {
	c, err := net.Dial("tcp", ":37374")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	conn := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(c, jsonrpc2.VSCodeObjectCodec{}), &clientHandler{})

	fullpath, err := filepath.Abs(".")
	uri := util.PathToURI(fullpath)
	if err != nil {
		panic(err)
	}

	fmt.Println(uri)
	params := InitializeParams{
		RootPath: fullpath,
		RootURI:  uri,
	}

	var got interface{}
	err = conn.Call(context.Background(), "initialize", params, &got)

	b, err := json.Marshal(got)

	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(string(b))

	fmt.Println("initial Call done!! next handshake...")

	// next, send handShake method...
	time.Sleep(time.Second * 1)
	handshakeReq := handshakeRequest{}
	handshakeRes := handshakeResponse{}
	err = conn.Call(context.Background(), "gopls/handshake", handshakeReq, &handshakeRes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("handshake result: %+v\n", handshakeRes)

	// wait for message from server...
	ctx, cancel := context.WithCancel(context.Background())
	go waitReq(ctx)

	time.Sleep(time.Second * 1)

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

	fmt.Println("cancel done!! program exit...")
}

func waitReq(ctx context.Context) {
	for {
		time.Sleep(time.Second * 1)
		select {
		case <-ctx.Done():
			fmt.Printf("Done\n")
			return
		default:
		}
	}
}
