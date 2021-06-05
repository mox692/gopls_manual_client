package client

import (
	"context"
	"log"
	"os"

	"github.com/sourcegraph/jsonrpc2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	logfile string
	port    string
	logger  *log.Logger
}

type inputConfig struct {
	Logfile string `yaml:"logfile"`
	Port    string `yaml:"port"`
}

func LoadConfig(yamlfile string) (*Config, error) {

	f, err := os.Open(yamlfile)
	if err != nil {
		log.Println("efm-langserver: no configuration file")
		// return default configs.
		return &Config{logger: log.New(os.Stdout, "", log.LstdFlags)}, nil
	}
	defer f.Close()

	var config = inputConfig{}
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return initConfig(&config)
}

func initConfig(config *inputConfig) (*Config, error) {
	if logfile := config.Logfile; logfile != "" {
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return &Config{
			logfile: config.Logfile,
			logger:  log.New(f, "", log.LstdFlags),
			port:    config.Port,
		}, nil
	}

	return &Config{
		logfile: config.Logfile,
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		port:    config.Port,
	}, nil
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
