package client

import (
	"context"

	"github.com/mox692/gopls_manual_client/protocol"
)

func (req *ClientRequester) didChange() error {
	var didChangeResult interface{}
	// TODO: ここのparamsは別所から取ってくるようにする

	err := req.conn.Call(context.Background(), "textDocument/didChange", didChangeParams, &didChangeResult)
	if err != nil {
		return err
	}
	return nil
}

var didChangeParams = protocol.DidChangeTextDocumentParams{
	TextDocument: protocol.VersionedTextDocumentIdentifier{
		Version: 2,
		TextDocumentIdentifier: protocol.TextDocumentIdentifier{
			URI: "dummyURI",
		},
	},
}
