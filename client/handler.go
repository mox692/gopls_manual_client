package client

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

type ClientHandler struct {
}

func (c *ClientHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "textDocument/didChange":
		c.handleTextDocumentDidChange()
		// TODO: impl
	}

}
