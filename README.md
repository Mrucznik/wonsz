# Wonsz

---

**W**rapper **O**f **N**aughty **S**nake**Z**

---

**The best of Viper & Cobra combined.**  
Ready to go solution for configurable CLI programs.

## What is this?

You get configuration struct, that fields are automatically bound to:
1. configuration file
2. environment variables
3. commandline flags

## Why?

- Let's say you want to write configurable app.  
> So, you use viper to load configuration from file.
You fetch your configuration fields by `viper.Get(key)`.  
- But it sucks to not have autocompletion from IDE.  
> So you marshall your config to struct.  
- But let's say, you dockerized your app, and when you run containers, you want also to manage config by environment variables.  
> So you use AutomaticEnv to get env variables. 
- But it marshall to struct only when you bind specific environment variable by name.  
> So you bind them. 
- But you have config struct field named like: ThisIsMyConfigField, so you must set THISISMYCONFIGFIELD env variable, which is not really readable and nice.
> You could write a wrapper.
- And let's say, you also want to run you app like an CLI app.  
> So you use cobra.   
- But you want to overwrite some configuration fields by the command line flags.  
> So you use viper.BindPFlag to bind some flags to your config structs.  
- But now you have  in your code for the same config field and pretty complicated initialization logic.
- And you end up with 3 different names of the same config field, and pretty complicated initialization logic. Also you must remember to add proper code when adding new field to configuration.

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


### Create a config struct in your code
```go
// config/config.go
package config

import (
	"github.com/Mrucznik/wonsz"
)

type Config struct {
	wonsz.Config
	
	// Here we declare fields with additional configuration tags
	SampleConfigField string `mapstructure:"sample_config_field" default:"default value"`
}

var cfg *Config

func Get() *Config {
	return cfg
}

func init() {
	wonsz.InitializeConfig(cfg, wonsz.ConfigOpts{})
}
```

### Initialize config in your cobra application

```go
package config

import "github.com/Mrucznik/wonsz"

```

### Use configuration fields in your code
```go
package main

import "yourapp/config"

var cfg = config.Get()

func main() {
	fmt.Println(cfg.SampleConfigField)	
}
```

### Configure and run your application with:
- **configuration files**
  - *config.json*
    ```json
    {
      "sample_config_field": "some value" 
    }
    ```
  - *config.yaml*
    ```yaml
    sample_config_field: "some value"
    ``` 
  - *config.toml*
    ```toml
    sample_config_field = "some value"
    ```
- **environment variables**
  ```shell
  SAMPLE_CONFIG_FIELD="some value" go run main.go
  ```
- or **command-line flags**
  ```shell
  go run main.go --sample-config-field="some value"
  ``` 

## Detailed configuration options

TODO.