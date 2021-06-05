package client

import (
	"context"
	"log"

	"github.com/sourcegraph/jsonrpc2"
)

type Config struct {
	logger *log.Logger
}

func LoadConfig() *Config {
	return &Config{}
}

type ClientHandler struct {
	logger *log.Logger
}

func NewHandler(config *Config) *ClientHandler {
	handler := &ClientHandler{
		logger: config.logger,
	}
	return handler
}

func (c *ClientHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "textDocument/didChange":
		c.handleTextDocumentDidChange()
		// TODO: impl
	}

}
