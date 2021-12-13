package wonsz

import (
	"github.com/spf13/cobra"
	"testing"
)

type TestConfig struct {
	AlwaysInConfig string `sometag:"xd" boolTag:""`
	TestField      string `name:"TestString" booltag:"" kek01:"pur"`
}

func TestInitializeConfig(t *testing.T) {
	err := Wonsz(TestConfig{}, &cobra.Command{}, ConfigOpts{})
	if err != nil {
		t.Fatal(err)
	}
}
