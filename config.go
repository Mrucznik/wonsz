package main

import (
	"fmt"
	"github.com/sevlyar/retag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
)

var cfgOpts ConfigOpts
var cfg interface{}

type ConfigOpts struct {
	Prefix      string
	ConfigPaths []string
	ConfigType  string
	ConfigName  string
}

func InitializeViper() {
	viper.SetEnvPrefix(cfgOpts.Prefix)

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
}

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

func Wonsz(config interface{}, rootCmd *cobra.Command, options ConfigOpts) error {
	cfgOpts = options
	cfg = retag.Convert(config, MapstructureRetagger{})

	cobra.OnInitialize(InitializeViper)

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
				logrus.Error("arrays not supported yet.")
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
				logrus.Error("arrays not supported yet.")
			default:
				continue
			}
		}

		targetFlag := flags.Lookup(dashedName)
		if targetFlag == nil {
			return fmt.Errorf("nie ma flag")
		}
		err := viper.BindPFlag(underscoredName, targetFlag)
		if err != nil {
			return err
		}
	}
	return nil
}
