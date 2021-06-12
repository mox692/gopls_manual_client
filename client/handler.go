package client

import (
	"context"
	"log"
	"os"

	"github.com/sourcegraph/jsonrpc2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Logfile *os.File
	Port    string
	Logger  *log.Logger
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
		return &Config{Logger: log.New(os.Stdout, "", log.LstdFlags)}, nil
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
		log.Printf("logging to : %s\n", logfile)
		return &Config{
			Logfile: f,
			Logger:  log.New(f, "", log.LstdFlags),
			Port:    config.Port,
		}, nil
	}
	return &Config{
		Logfile: os.Stdout,
		Logger:  log.New(os.Stdout, "", log.LstdFlags),
		Port:    config.Port,
	}, nil
}

type ClientHandler struct {
	logger *log.Logger
}

func NewHandler(config *Config) *ClientHandler {
	handler := &ClientHandler{
		logger: config.Logger,
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
