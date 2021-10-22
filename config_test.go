package main

import (
	"github.com/spf13/cobra"
	"testing"
)

type TestConfig struct {
	AlwaysInConfig string `sometag:"xd" boolTag:""`
	TestField      string `name:"TestString" booltag:"" kek01:"pur"`
}

func TestInitializeConfig(t *testing.T) {
	Wonsz(TestConfig{}, &cobra.Command{}, ConfigOpts{})
}
