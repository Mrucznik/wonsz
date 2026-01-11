package tests

import (
	"testing"

	"github.com/Mrucznik/wonsz"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestConfigCobraFlags(t *testing.T) {
	var config Config

	cmd := &cobra.Command{
		Use:           "test",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	opts := wonsz.ConfigOpts{
		EnvPrefix:             "WONSZ",
		ConfigPaths:           []string{"."},
		ConfigType:            "json",
		ConfigName:            "empty_config",
		Viper:                 viper.New(),
		IgnoreViperBindErrors: false,
	}

	if err := wonsz.BindConfig(&config, cmd, opts); err != nil {
		t.Fatal(err)
	}

	cmd.SetArgs([]string{
		"--app-name", "ExampleApp",
		"--version", "1.0.0",
		"--debug=true",
		"--start-time", "2025-01-01T12:30:45Z",
		"--start-duration", "1h30m",

		"--server-host", "0.0.0.0",
		"--server-port", "8080",
		"--server-timeouts-read", "5s",
		"--server-timeouts-write", "10s",

		"--database-driver", "postgres",
		"--database-max-connections", "20",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if err := assertConfigWithoutStructArrays(config); err != nil {
		t.Fatal(err)
	}
}
