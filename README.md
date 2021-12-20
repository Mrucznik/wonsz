![img](wonsz.png)

---

**W**rapper **O**f **N**aughty **S**nake**Z**

---

**The best of Viper & Cobra combined.**  
Ready to go solution for configurable CLI programs.

## What it does?

It creates configuration struct, that fields are automatically bound to:

1. configuration file
2. environment variables
3. commandline flags

## Why?

- Let's say you want to write configurable app.

> So, you use viper to load configuration from file. You fetch your configuration fields by `viper.Get(key)`.

- But it sucks to not have autocompletion from IDE.

> So you marshall your config to struct.

- But let's say, you dockerized your app, and when you run containers, you want also to manage config by environment
  variables.

> So you use AutomaticEnv to get env variables.

- But it marshall to struct only when you bind specific environment variable by name.

> So you bind them.

- But you have config struct field named like: ThisIsMyConfigField, so you must set THISISMYCONFIGFIELD env variable,
  which is not really readable and nice.

> You could write a wrapper.

- And let's say, you also want to run you app like an CLI app.

> So you use cobra.

- But you want to overwrite some configuration fields by the command line flags.

> So you use viper.BindPFlag to bind some flags to your config structs.

- But now you have in your code for the same config field and pretty complicated initialization logic.
- And you end up with 3 different names of the same config field, and pretty complicated initialization logic. Also you
  must remember to add proper code when adding new field to configuration.

> So you use this library, and then you just **create 1 config struct** without any tags, initialize it,
> and you have **all 3 ways of configuring you app** (by configuration file, by environment variables and by command flags) out of the box and in one place.  
> And you have all the above problems resolved.

- Awesome!

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

var config *Configuration

type Configuration struct {
	// Here we declare configuration fields. No need to add any tags.
	SnakeName string
}

func main() {
	wonsz.Wonsz(config, &cobra.Command{Run: execute}, wonsz.ConfigOpts{})
}

func execute(_ *cobra.Command, _ []string) {
	fmt.Printf("Application config: %+v\n", config)
}
```

This is the simplest example, more detailed [you will find here](example/example.go).

### Integrate with existing application

1. Create file config/config.go
    ```go
    // config.go file
    package config
   
   // TODO: more code
    ```
2. xd
3. ???
4. profit

### Configure and run your application with

- **default struct values**
  ```cgo
  config := &Config{
      SnakeName: "nope-rope",
  }
  ```
- **configuration files**
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