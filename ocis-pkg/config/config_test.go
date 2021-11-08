package config

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	yBytes, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(yBytes))
}
