package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/sourcegraph/go-langserver/langserver/util"
	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
	// "golang.org/x/tools/internal/lsp/protocol"
)

/**
 * these struct is referenced from  golang/tools/internal/lsp/protocol/tsprotocol.go
 * ref: https://github.com/golang/tools/blob/master/internal/lsp/protocol/tsprotocol.go
 */
type clientHandler struct {
}

type InitializeParams struct {
	RootPath string
	RootURI  string
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

func main() {
	c, err := net.Dial("tcp", ":1234")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

	conn := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(c, jsonrpc2.VSCodeObjectCodec{}), &clientHandler{})

	fmt.Printf("%+v\n", conn)
	done := make(chan bool)

	fullpath, err := filepath.Abs(".")
	uri := util.PathToURI(fullpath)
	if err != nil {
		panic(err)
	}

	fmt.Println(uri)
	params := InitializeParams{
		RootPath: fullpath,
	}

	var got interface{}
	err = conn.Call(context.Background(), "initialize", params, &got)

	// fmt.Printf("%+v\n", got)

	b, err := json.Marshal(got)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(string(b))

	fmt.Println("initial Call done!!!")

	var got2 interface{}
	params2 := &DocumentHighlightParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///Users/kimuramotoyuki/go/src/github.com/mox692/gopls_manual_client",
			},
			Position: Position{
				Line: 3,
			},
		},
	}

	err = conn.Call(context.Background(), "textDocument/documentHighlight", params2, &got2)
	if err != nil {
		panic(err)
	}

	b2, _ := json.Marshal(got2)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(string(b2))

	fmt.Println("b2 Call done!!!")

	go func() {
		<-conn.DisconnectNotify()
		done <- true
	}()
}
