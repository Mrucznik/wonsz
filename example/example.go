package main

import (
	"fmt"

	"github.com/Mrucznik/wonsz"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Config is our example config struct.
type Config struct {
	// Here you can define configuration fields that will be automatically bound to a config file, env and flags.
	// You don't need to specify additional field tags.
	// Fields in config will be underscore-separated (snake case).
	// Environment variables will be all uppercase-underscore-separated (screaming snake case).
	// Flags will be dash-separated (kebab-case).
	SnakeName   string // For example, this will become: snake_name in the file, SNAKE_NAME for env, and snake-name for a flag.
	SnakeLength int
	SnakeHappy  bool
}

var config *Config

// Here we define our root command with cobra.
// If you don't want to use Wonsz with cobra CLI, just pass nil as rootCmd.
var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "A brief description of your application.",
	Long:  `Detailed description of your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Here's where you define what will happen when you run your application.
		fmt.Printf("Application config: %+v\n", config)
	},
}

// You should init your config as first what your application will do,
// to dodge some problems with referring to the config field that is uninitialized.
func init() {
	// We create a struct that will be used for binding values to.
	config = &Config{
		// Here you can specify some default values when there is no configuration available.
		SnakeName: "nope-rope",
	}
	err := wonsz.BindConfig(config, rootCmd,
		// You can specify some additional options to Wonsz, so configuration settings can better meet your needs.
		// If you want to go with default, you can pass the empty struct here.
		wonsz.ConfigOpts{
			EnvPrefix:   "WONSZ",
			ConfigPaths: []string{"."},
			ConfigType:  "json",
			ConfigName:  "example",
		},
	)
	// Here we should check, if config was created successfully.
	if err != nil {
		logrus.WithError(err).Fatal("cannot load config")
		return
	}
}

func main() {
	// we execute our application
	err := rootCmd.Execute()
	if err != nil {
		logrus.WithError(err).Fatal("application failed to execute")
	}
}

// You can run your application and overwrite any config field by ENVIRONMENT_VARIABLE:
// $ WONSZ_SNAKE_NAME="caution ramen" go run .
// Or command line flag:
// $ go run . --snake-name "danger noodle"
// Or without any modifications (then we will get config values from the example.json file):
// $ go run .

// You can also get some help with config variables available to override:
// $ go run . --help
