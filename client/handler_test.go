package client

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testYamlFile := "/Users/kimuramotoyuki/go/src/github.com/mox692/gopls_manual_client/config_test.yaml"
	expected := Config{
		Port: "37375",
	}
	config, err := LoadConfig(testYamlFile)
	if err != nil {
		t.Errorf("Fail: %+v", err)
	}
	if config.Port != expected.Port {
		t.Errorf("port not match. expect %s, but got %s\n", expected.Port, config.Port)
	}
}
