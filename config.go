package wonsz

import (
	"fmt"
	"github.com/sevlyar/retag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
)

var cfgOpts *ConfigOpts
var cfg interface{}

// ConfigOpts provide additional options to configure Wonsz.
type ConfigOpts struct {
	// Environment variables prefix.
	// E.g. if your prefix is "wonsz", the env registry will look for env variables that start with "WONSZ_".
	EnvPrefix string

	// Paths to search for the config file in.
	ConfigPaths []string

	// Type of the configuration file, e.g. "json".
	ConfigType string

	// Name for the config file. Does not include extension.
	ConfigName string
}

// TODO: Here some description
func Wonsz(config interface{}, rootCmd *cobra.Command, options ConfigOpts) error {
	cfgOpts = &options
	cfg = retag.Convert(config, mapstructureRetagger{})

	cobra.OnInitialize(initializeViper)

	confType := reflect.TypeOf(cfg).Elem()
	for i := 0; i < confType.NumField(); i++ {
		field := confType.Field(i)
		if field.Anonymous {
			continue
		}
		dashedName := CamelCaseToDashedLowered(field.Name)
		underscoredName := CamelCaseToUnderscoredLowered(field.Name)

		flags := rootCmd.PersistentFlags()
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
				return fmt.Errorf("arrays not supported yet")
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
				return fmt.Errorf("arrays not supported yet")
			default:
				continue
			}
		}

		targetFlag := flags.Lookup(dashedName)
		if targetFlag == nil {
			// TODO: better errors
			return fmt.Errorf("nie ma flag")
		}
		err := viper.BindPFlag(underscoredName, targetFlag)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: here some description
func initializeViper() {
	viper.SetEnvPrefix(cfgOpts.EnvPrefix)

	for _, path := range cfgOpts.ConfigPaths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigType(cfgOpts.ConfigType)
	viper.SetConfigName(cfgOpts.ConfigName)

	bindEnvsAndSetDefaults()
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Infof("Config file not found.")
	} else {
		logrus.Infof("Using config file: %v.", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.WithError(err).Fatal("Cannot unmarshall config into Config struct.")
	}

	// TODO: do we relly don't want to store this configuration globally and want to garbage collector to remove it?
	cfgOpts = nil
}

// TODO: here some description
func bindEnvsAndSetDefaults() {
	el := reflect.TypeOf(cfg).Elem()
	for i := 0; i < el.NumField(); i++ {
		field := el.Field(i)
		defaultVal := field.Tag.Get("default")
		mapping := field.Tag.Get("mapstructure")

		if defaultVal != "" {
			viper.SetDefault(mapping, defaultVal)
		} else {
			err := viper.BindEnv(mapping)
			if err != nil {
				logrus.Fatal(err)
			}
		}
	}
}
