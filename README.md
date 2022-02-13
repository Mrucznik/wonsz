![img](wonsz.png)

---

**W**rapper **O**f **N**aughty **S**nake**Z**

---

**The best of Viper & Cobra combined.**  
Ready to go solution for configurable CLI programs.

![example workflow](https://github.com/Mrucznik/wonsz/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Mrucznik/wonsz)](https://goreportcard.com/report/github.com/Mrucznik/wonsz)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/Mrucznik/wonsz)](https://pkg.go.dev/mod/github.com/Mrucznik/wonsz)

## What does it do?

It creates a configuration struct, that fields are automatically bound to:

1. configuration file
2. environment variables
3. command line flags

## Why?

- Let's say you want to write a configurable app.

> So, I use viper to load configuration from file. I fetch configuration fields by `viper.Get(key)`.

- But it sucks to not have autocompletion from IDE.

> So, I marshall your config to a struct.

- But let's say, you dockerized your app and when you run containers, you want also to manage config by environment
  variables.

> So, I use AutomaticEnv to get env variables.

- But it marshalls to struct only when you bind specific environment variables by name.

> I would bind them by viper.BindEnv().

- But you have config struct field named like: ThisIsMyConfigField, so you must set THISISMYCONFIGFIELD env variable,
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

Import dependency into your project.

```shell
go get -d github.com/Mrucznik/wonsz
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
	wonsz.BindConfig(&config, // pointer to the configuration struct
		rootCmd,            // root cobra command
		wonsz.ConfigOpts{}) // Wonsz configuration options

	rootCmd.Execute()
}

func execute(_ *cobra.Command, _ []string) {
	fmt.Printf("Application config: %+v\n", config)
}
```

This is the simplest example, more detailed [you will find here](example/example.go).

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
	wonsz.BindConfig(&config.Config, rootCmd, wonsz.ConfigOpts{})
	
	// other code
}

```
3. Done!

### Configure and run your application with

- **default struct values**
  ```cgo
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
  go run main.go --snake_name="danger noodle"
  ``` 

## Detailed configuration options

You can find more information by checking out [example app](example/example.go).
