![img](wonsz.png)

---

**W**rapper **O**f **N**aughty **S**nake**Z**

---

**The best of Viper & Cobra combined.**  
Bind your config struct to CLI flags, environment variables, and configuration files.

![example workflow](https://github.com/Mrucznik/wonsz/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Mrucznik/wonsz)](https://goreportcard.com/report/github.com/Mrucznik/wonsz)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/Mrucznik/wonsz)](https://pkg.go.dev/mod/github.com/Mrucznik/wonsz)

## What does it do?

It creates a configuration struct that fields are automatically bound to:

1. configuration file
2. environment variables
3. command line flags

## Why?

- Let's say you want to write a configurable app.

> So, I use viper to load configuration from a file. I fetch configuration fields by `viper.Get(key)`.

- But it sucks to not have autocompletion from the IDE.

> So, I marshall your config to a struct.

- But let's say, you dockerized your app, and when you run containers, you want also to manage config by environment
  variables.

> So, I use AutomaticEnv to get env variables.

- But it marshalls to struct only when you bind specific environment variables by name.

> I would bind them by viper.BindEnv().

- But you have a config struct field named like: ThisIsMyConfigField, so you must set THISISMYCONFIGFIELD env variable,
  which is not really readable and nice.

> :/

- And let's say, you also want to run your app like a CLI app.

> I would use cobra.

- But you may also want configuration fields to be overwritten with values from the command line flags.

> So, I use viper.BindPFlag to bind some flags to your config structs.

- And you end up with 3 different names of the same config field and pretty complicated initialization logic. Also, you
  must remember to add proper code when adding a new field to the configuration, so every way of loading the config field is properly handled.

> So I use this library, and then you just **create 1 config struct** without any tags, initialize it,
> and you have **all 3 ways of configuring your app** (by the configuration file, by environment variables, and by command flags) out of the box and in one place.  
> And I have all the above problems resolved!

- Awessssome!

## How to install?

Requires Go 1.25 or newer. Import dependency into your project.

```shell
go get github.com/Mrucznik/wonsz
```

## How to use?

### Simplest application

```go
// main.go file
package main

import (
	"fmt"
	"github.com/Mrucznik/wonsz"
	"github.com/spf13/cobra"
)

var config Configuration

type Configuration struct {
	// Here we declare configuration fields. No need to add any tags.
	SnakeName string
}

var rootCmd = &cobra.Command{Run: execute}

func main() {
	err := wonsz.BindConfig(&config, // pointer to the configuration struct
		rootCmd,            // root cobra command
		wonsz.ConfigOpts{}) // Wonsz configuration options
	if err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func execute(_ *cobra.Command, _ []string) {
	fmt.Printf("Application config: %+v\n", config)
}
```

This is the simplest example, more production-ready and detailed [you will find here](example/example.go). You can also look at [tests](tests/config_test.go).

### Integrate with an existing application

1. Create file config/config.go
    ```go
    // config.go file
    package config
   
    var Config Configuration

    type Configuration struct {
    // Here we declare configuration fields. No need to add any tags.
        SnakeName string
    }
    ```
2. Bind created config structure to cobra & viper using wonsz.BindConfig()
```go
// cmd/root.go
package cmd

import (
    "github.com/Mrucznik/wonsz"
	"github.com/You/your-project/config"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	// your root cobra command
}

func init() {
	if err := wonsz.BindConfig(&config.Config, rootCmd, wonsz.ConfigOpts{}); err != nil {
		panic(err)
	}

	// other code
}

```
3. Done!

### Configure and run your application with

- **default struct values**
  ```go
  config := &Config{
      SnakeName: "nope-rope",
  }
  ```
- **configuration files** e.g.
    - *config.json*
      ```json
      {
        "snake_name": "hazard spaghetti" 
      }
      ```
    - *config.yaml*
      ```yaml
      snake_name: "judgemental shoelace"
      ``` 
    - *config.toml*
      ```toml
      snake_name = "slippery tube dude"
      ```
- **environment variables**
  ```shell
  SNAKE_NAME="caution ramen" go run main.go
  ```
- **command-line flags**
  ```shell
  go run main.go --snake-name="danger noodle"
  ``` 

## Key concepts

- keys in configuration files should be in snake_case
- environment variables are in SCREAMING_SNAKE_CASE
- command-line flags names should be dash-separated

## Supported field types

Strings, booleans, all int/uint/float variants, `time.Duration`, `time.Time` (RFC 3339),
`net.IP`, `net.IPNet`, slices of strings/bools/ints/floats/`time.Duration`/`net.IP`,
string arrays, and `map[string]string`/`map[string]int`/`map[string]int64`.

Nested structs (also behind pointers) are fully supported — their fields get prefixed
names, e.g. `Server.Port` becomes `server.port` in the file, `SERVER_PORT` in env,
and `--server-port` on the command line.

## Field tags

All tags are optional:

| Tag | Effect |
|-----|--------|
| `mapstructure:"custom_name"` | Overrides the derived name for all bindings (flag becomes `--custom-name`). |
| `default:"value"` | Default value used when no source provides one. |
| `usage:"help text"` | Usage string shown in `--help` for the flag. |
| `shortcut:"p"` | Single-character flag shorthand, e.g. `-p`. |
| `wonsz:"flag-ignore"` | Skips binding the field to a command-line flag (env and file still work). |
| `wonsz:"-"` | Excludes the field from all bindings entirely. |

## Configuration options

`ConfigOpts` fields:

- `EnvPrefix` — prefix for environment variables (`"WONSZ"` → `WONSZ_SNAKE_NAME`).
- `ConfigPaths`, `ConfigType`, `ConfigName` — where and how to look for the config file.
- `Viper` — pass your own viper instance (defaults to the global one).
- `IgnoreFlagBindErrors` — skip fields that cannot be bound to flags instead of returning an error.
- `WatchConfig` — watch the config file and re-unmarshal the struct when it changes.

## Typed instances

`wonsz.New` returns a typed `*wonsz.Wonsz[T]` instance with `Get() *T` (no type
assertions) and `Viper()`. Instances are independent, so several configs can
coexist:

```go
w, err := wonsz.New(&cfg, rootCmd, wonsz.ConfigOpts{Viper: viper.New()})
name := w.Get().SnakeName
```

`wonsz.BindConfig` is a convenience wrapper for when you keep your own pointer
to the struct and only need the error.

## Cobra integration notes

Wonsz chains config initialization into the root command's
`PersistentPreRunE`; your own hook (if any) runs right after it. If a
subcommand defines its own `PersistentPreRun(E)`, cobra skips the root hook —
set `cobra.EnableTraverseRunHooks = true` to run both.

## Stability

Wonsz follows semantic versioning. Starting with v1.0.0 the public API is
stable: breaking changes only happen in a new major version.

## More examples

You can find more information by checking out [example app](example/example.go).
