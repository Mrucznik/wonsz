package tests

import (
	"testing"

	"github.com/Mrucznik/wonsz"
)

func TestJsonConfig(t *testing.T) {
	var config Config

	err := wonsz.BindConfig(&config, nil, wonsz.ConfigOpts{
		ConfigName:  "json_config",
		ConfigPaths: []string{"."},
		ConfigType:  "json",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = assertConfig(config); err != nil {
		t.Fatal(err)
	}
}
