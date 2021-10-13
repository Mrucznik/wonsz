package wonsz

import (
	"testing"
)

type TestConfig struct {
	Config

	TestField string `name:"TestString" booltag:"" kek01:"pur"`
}

func TestInitializeConfig(t *testing.T) {
	InitializeConfig(TestConfig{}, nil)
}
