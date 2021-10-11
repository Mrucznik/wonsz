# Wonsz

---

**W**rapper **O**f **N**aughty **S**nakes **Z**oo

---

**The best of Viper & Cobra combined.**  
Ready to go solution for configurable CLI programs.

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
	wonsz.InitializeConfig(cfg)
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