package tests

import (
	"testing"

	"github.com/Mrucznik/wonsz"
	"github.com/spf13/viper"
)

func TestConfigEnv(t *testing.T) {
	t.Setenv("WONSZ_APP_NAME", "ExampleApp")
	t.Setenv("WONSZ_VERSION", "1.0.0")
	t.Setenv("WONSZ_DEBUG", "true")
	t.Setenv("WONSZ_START_TIME", "2025-01-01T12:30:45Z")
	t.Setenv("WONSZ_START_DURATION", "1h30m")

	t.Setenv("WONSZ_SERVER_HOST", "0.0.0.0")
	t.Setenv("WONSZ_SERVER_PORT", "8080")
	t.Setenv("WONSZ_SERVER_TIMEOUTS_READ", "5s")
	t.Setenv("WONSZ_SERVER_TIMEOUTS_WRITE", "10s")

	t.Setenv("WONSZ_DATABASE_DRIVER", "postgres")
	t.Setenv("WONSZ_DATABASE_MAX_CONNECTIONS", "20")

	var config Config
	opts := wonsz.ConfigOpts{
		EnvPrefix:   "WONSZ",
		ConfigPaths: []string{"."},
		ConfigType:  "json",
		ConfigName:  "empty_config",
		Viper:       viper.New(),
	}

	if err := wonsz.BindConfig(&config, nil, opts); err != nil {
		t.Fatal(err)
	}

	if err := assertConfigWithoutStructArrays(config); err != nil {
		t.Fatal(err)
	}
}
