package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Config struct {
	TestString string `mapstructure:"test_string"`
	TestInt    int    `mapstructure:"test_int"`
	TestBool   bool
}

var config *Config

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "A brief description of your application",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Application config: %+v\n", config)
	},
}

func init() {
	config = &Config{}
	Wonsz(config, rootCmd, ConfigOpts{
		Prefix:      "WONSZ",
		ConfigPaths: []string{"."},
		ConfigType:  "toml",
		ConfigName:  "wonsz",
	})

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.main.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cobra.CheckErr(cobra.MarkFlagRequired(rootCmd.Flags(), "toggle"))
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
