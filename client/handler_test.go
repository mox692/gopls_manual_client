package client

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testYamlFile := "/Users/kimuramotoyuki/go/src/github.com/mox692/gopls_manual_client/config_test.yaml"
	expected := Config{
		logfile: "/Users/kimuramotoyuki/go/src/github.com/mox692/gopls_manual_client/client_debuglog",
		port:    "37375",
	}
	config, err := LoadConfig(testYamlFile)
	if err != nil {
		t.Errorf("Fail: %+v", err)
	}
	if config.logfile != expected.logfile {
		t.Errorf("logfile not match. expect %s, but got %s\n", expected.logfile, config.logfile)
	}
	if config.port != expected.port {
		t.Errorf("port not match. expect %s, but got %s\n", expected.port, config.port)
	}
}
