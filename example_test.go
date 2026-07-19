package wonsz_test

import (
	"fmt"
	"os"

	"github.com/Mrucznik/wonsz"
	"github.com/spf13/viper"
)

func ExampleNew() {
	_ = os.Setenv("SNAKE_NAME", "danger noodle")
	defer func() { _ = os.Unsetenv("SNAKE_NAME") }()

	type Config struct {
		SnakeName   string
		SnakeLength int `default:"5"`
	}

	var cfg Config
	w, err := wonsz.New(&cfg, nil, wonsz.ConfigOpts{Viper: viper.New()})
	if err != nil {
		panic(err)
	}

	// Get returns the typed pointer, no type assertion needed.
	fmt.Println(w.Get().SnakeName, w.Get().SnakeLength)
	// Output: danger noodle 5
}

func ExampleNew_tags() {
	_ = os.Setenv("CUSTOM_NAME", "hazard spaghetti")
	defer func() { _ = os.Unsetenv("CUSTOM_NAME") }()

	type Config struct {
		Renamed  string `mapstructure:"custom_name"`
		Excluded string `wonsz:"-"`
	}

	var cfg Config
	if err := wonsz.BindConfig(&cfg, nil, wonsz.ConfigOpts{Viper: viper.New()}); err != nil {
		panic(err)
	}

	fmt.Printf("%q %q\n", cfg.Renamed, cfg.Excluded)
	// Output: "hazard spaghetti" ""
}
