package wonsz

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
)

type IConfig interface {
}

type Config struct {
	AlwaysInConfig string `sometag:"xd" boolTag:""`
}

type ConfigOpts struct {
	Prefix string
}

func InitializeConfig(config IConfig, _ ...ConfigOpts) error {
	// TODO
	// 1. get all field names
	// 2. get all tags for fields
	// 3. bind env variables - XXX_XXX_XXX
	// 4. create flags --xxx-xxx-xxx
	// 5. load config

	confType := reflect.TypeOf(config)
	fmt.Println("type: ", confType)

	// get all field names & tags
	for i := 0; i < confType.NumField(); i++ {
		field := confType.Field(i)
		if field.Anonymous {
			continue
		}
		//camelCaseName := field.Name
		dashedName, underscoredName := ConvertFromCamelCase(field.Name)

		// 3. bind env variables - XXX_XXX_XXX
		if defaultVal, ok := field.Tag.Lookup("default"); ok {
			viper.SetDefault(underscoredName, defaultVal)
		} else {
			err := viper.BindEnv(underscoredName)
			if err != nil {
				// TODO: better error
				return err
			}
		}

		// 4. create flags --xxx-xxx-xxx
		rootCommand := cobra.Command{}
		flags := rootCommand.PersistentFlags()
		usageHint := field.Tag.Get("usage")
		if shortcut, ok := field.Tag.Lookup("shortcut"); ok {
			switch field.Type.Kind() {
			case reflect.String:
				flags.StringP(dashedName, shortcut, "", usageHint)
			case reflect.Int64:
				flags.Int64P(dashedName, shortcut, 0, usageHint)
			case reflect.Int32:
				flags.Int32P(dashedName, shortcut, 0, usageHint)
			case reflect.Int16:
				flags.Int16P(dashedName, shortcut, 0, usageHint)
			case reflect.Int8:
				flags.Int8P(dashedName, shortcut, 0, usageHint)
			case reflect.Int:
				flags.IntP(dashedName, shortcut, 0, usageHint)
			case reflect.Float64:
				flags.Float64P(dashedName, shortcut, 0, usageHint)
			case reflect.Float32:
				flags.Float32P(dashedName, shortcut, 0, usageHint)
			case reflect.Bool:
				flags.BoolP(dashedName, shortcut, false, usageHint)
			case reflect.Array, reflect.Slice, reflect.Map:
				//TODO: arrays
				fmt.Println("arrays not supported yet.")
			default:
				continue
			}
		} else {
			switch field.Type.Kind() {
			case reflect.String:
				flags.String(dashedName, "", usageHint)
			case reflect.Int64:
				flags.Int64(dashedName, 0, usageHint)
			case reflect.Int32:
				flags.Int32(dashedName, 0, usageHint)
			case reflect.Int16:
				flags.Int16(dashedName, 0, usageHint)
			case reflect.Int8:
				flags.Int8(dashedName, 0, usageHint)
			case reflect.Int:
				flags.Int(dashedName, 0, usageHint)
			case reflect.Float64:
				flags.Float64(dashedName, 0, usageHint)
			case reflect.Float32:
				flags.Float32(dashedName, 0, usageHint)
			case reflect.Bool:
				flags.Bool(dashedName, false, usageHint)
			case reflect.Array, reflect.Slice, reflect.Map:
				//TODO: arrays
				fmt.Println("arrays not supported yet.")
			default:
				continue
			}
		}

		err := viper.BindPFlag(underscoredName, flags.Lookup(dashedName))
		if err != nil {
			// TODO: better error
			return err
		}

		// get all tags
		//for j, key := range GetTagsForField(field) {
		//	fmt.Printf("tag %d: %+v\n", j, key)
		//}
	}

	// 5. load config
	err := viper.Unmarshal(&config)
	if err != nil {
		// TODO: better error
		return err
	}

	return nil
}
