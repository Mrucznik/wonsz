package tests

import (
	"testing"

	"github.com/Mrucznik/wonsz"
)

func TestTomlConfig(t *testing.T) {
	var config Config

	err := wonsz.BindConfig(&config, nil, wonsz.ConfigOpts{
		ConfigName:  "toml_config",
		ConfigPaths: []string{"."},
		ConfigType:  "toml",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = assertConfig(config); err != nil {
		t.Fatal(err)
	}
}
