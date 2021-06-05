package client

import "github.com/sourcegraph/jsonrpc2"

type ClientRequester struct {
	conn   *jsonrpc2.Conn
	config *Config
}

func NewRequester(conn *jsonrpc2.Conn, config *Config) *ClientRequester {
	return &ClientRequester{
		conn:   conn,
		config: config,
	}
}

func (req *ClientRequester) CallMethod(method string, config *Config) error {
	switch method {
	case "textDocument/didChange":
		err := req.didChange()
		return err
	default:
		config.Logger.Printf("method request [%s] is not defined.\n", method)
	}
	return nil
}
