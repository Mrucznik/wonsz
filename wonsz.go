package wonsz

import (
	"fmt"
	"github.com/sevlyar/retag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"reflect"
)

// TODO: maybe ConfigOpts as variadic parameters?
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

// Get returns a config struct instance to which Wonsz binds configuration.
func Get() interface{} {
	return cfg
}

// TODO: next version: make an option to get viper and set lib to use own viper, not global instance
func GetViper() *viper.Viper {
	return viper.GetViper()
}

// TODO: next version: make an option to get command that flags will be binded to and set lib to use command
func GetCommand() *cobra.Command {
	return &cobra.Command{}
}

// BindConfig binds configuration structure to config file, environment variables and cobra command flags.
// The config parameter should be a pointer to configuration structure.
// You can pass nil to rootCmd, if you don't want to bind cobra command flags with config.
func BindConfig(config interface{}, rootCmd *cobra.Command, options ConfigOpts) error {
	if !(reflect.TypeOf(config).Kind() == reflect.Ptr && reflect.TypeOf(config).Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("config parameter is not a pointer to a structure. Maybe you should use & operator")
	}

	// prepare for processing
	cfgOpts = &options
	cfg = retag.Convert(config, mapstructureRetagger{})

	if rootCmd == nil { // only viper
		initializeViper()
		return nil
	}
	cobra.OnInitialize(initializeViper)

	confType := reflect.TypeOf(cfg).Elem()
	for i := 0; i < confType.NumField(); i++ {
		field := confType.Field(i)
		if field.Anonymous {
			continue
		}
		dashedName := camelCaseToDashedLowered(field.Name)
		underscoredName := camelCaseToUnderscoredLowered(field.Name)

		var err error
		flags := rootCmd.PersistentFlags()
		usageHint := field.Tag.Get("usage")
		if shortcut, ok := field.Tag.Lookup("shortcut"); ok {
			err = bindPFlag(flags, field, dashedName, shortcut, usageHint)
		} else {
			err = bindFlag(flags, field, dashedName, usageHint)
		}
		if err != nil {
			return err
		}

		targetFlag := flags.Lookup(dashedName)
		if targetFlag == nil {
			// TODO: better errors
			return fmt.Errorf("nie ma flag")
		}
		err = viper.BindPFlag(underscoredName, targetFlag)
		if err != nil {
			return err
		}
	}
	return nil
}

func bindPFlag(flags *pflag.FlagSet, field reflect.StructField, dashedName, shortcut, usageHint string) error {
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
	case reflect.Map:
		if field.Type.Elem().Kind() == reflect.String {
			flags.StringToStringP(dashedName, shortcut, map[string]string{}, usageHint)
		} else {
			return fmt.Errorf("unsupported flag %s type: %s. only map[string]string maps are supported", dashedName, field.Type.String())
		}
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.String {
			flags.StringSliceP(dashedName, shortcut, []string{}, usageHint)
		} else {
			return fmt.Errorf("unsupported flag %s type: %s. only string slices are supported", dashedName, field.Type.String())
		}
	case reflect.Array:
		if field.Type.Elem().Kind() == reflect.String {
			flags.StringArrayP(dashedName, shortcut, []string{}, usageHint)
		} else {
			return fmt.Errorf("unsupported flag %s type: %s. only string arrays are supported", dashedName, field.Type.String())
		}
	}
	return nil
}

func bindFlag(flags *pflag.FlagSet, field reflect.StructField, dashedName, usageHint string) error {
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
	case reflect.Array:
		if field.Type.Elem().Kind() == reflect.String {
			flags.StringArray(dashedName, []string{}, usageHint)
		} else {
			return fmt.Errorf("unsupported flag %s type: %s. only string arrays are supported", dashedName, field.Type.String())
		}
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.String {
			flags.StringSlice(dashedName, []string{}, usageHint)
		} else {
			return fmt.Errorf("unsupported flag %s type: %s. only string slices are supported", dashedName, field.Type.String())
		}
	case reflect.Map:
		if field.Type.Elem().Kind() == reflect.String {
			flags.StringToString(dashedName, map[string]string{}, usageHint)
		} else {
			return fmt.Errorf("unsupported flag %s type: %s. only map[string]string maps are supported", dashedName, field.Type.String())
		}
	}
	return nil
}

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
