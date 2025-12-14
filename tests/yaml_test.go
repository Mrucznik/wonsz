package tests

import (
	"testing"

	"github.com/Mrucznik/wonsz"
)

func TestYamlConfig(t *testing.T) {
	var config Config

	err := wonsz.BindConfig(&config, nil, wonsz.ConfigOpts{
		ConfigName:  "yaml_config",
		ConfigPaths: []string{"."},
		ConfigType:  "yaml",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = assertConfig(config); err != nil {
		t.Fatal(err)
	}
}
