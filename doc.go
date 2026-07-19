// Package wonsz binds a plain Go configuration struct to configuration files,
// environment variables and cobra command-line flags at once — the best of
// viper and cobra combined, without writing any binding boilerplate.
//
// Define a struct, call [New] (or the convenience wrapper [BindConfig]) and
// every field is readable from a config file (snake_case keys), environment
// variables (SCREAMING_SNAKE_CASE) and command-line flags (kebab-case), with
// the usual precedence: flags > env > config file > defaults.
//
//	type Config struct {
//		SnakeName string `default:"nope-rope" usage:"the snake's name"`
//	}
//
//	var cfg Config
//	w, err := wonsz.New(&cfg, rootCmd, wonsz.ConfigOpts{EnvPrefix: "APP"})
//
// Field behavior can be tuned with optional struct tags: mapstructure
// (custom key), default, usage, shortcut, wonsz:"flag-ignore" (no flag) and
// wonsz:"-" (excluded entirely). See ConfigOpts for file lookup, custom viper
// instances and config hot-reload via WatchConfig.
package wonsz
