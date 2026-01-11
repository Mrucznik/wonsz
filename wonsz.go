package wonsz

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/sevlyar/retag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	globalViper "github.com/spf13/viper"
)

var cfgOpts ConfigOpts
var cfg interface{}
var viper *globalViper.Viper

// ConfigOpts provide additional options to configure Wonsz.
type ConfigOpts struct {
	// Environment variables prefix.
	// E.g., if your prefix is "wonsz", the env registry will look for env variables that start with "WONSZ_".
	EnvPrefix string

	// Paths to search for the config file in.
	ConfigPaths []string

	// Type of the configuration file, e.g. "json".
	// Wonsz uses Viper for loading the configuration file,
	// so you can use any type of configuration file that Viper supports.
	ConfigType string

	// Name for the config file. Does not include extension.
	// If no configuration name is specified, Wonsz will not throw an error if the config file is not found.
	ConfigName string

	// Pass own viper instance. Default is a global viper instance.
	Viper *globalViper.Viper

	// If true, Wonsz will not return an error if a config field cannot be bound to a flag.
	IgnoreViperBindErrors bool
}

// Get returns a config struct instance to which Wonsz binds configuration.
func Get() interface{} {
	return cfg
}

// GetViper returns a viper instance used by Wonsz.
func GetViper() *globalViper.Viper {
	return viper
}

// BindConfig binds configuration structure to config file, environment variables and cobra command flags.
// The config parameter should be a pointer to the configuration structure.
// You can pass nil to rootCmd if you don't want to bind cobra command flags with config.
func BindConfig(config interface{}, rootCmd *cobra.Command, options ConfigOpts) error {
	if !(reflect.TypeOf(config).Kind() == reflect.Ptr && reflect.TypeOf(config).Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("config parameter is not a pointer to a structure. Maybe you should use & operator")
	}

	// prepare for processing
	cfgOpts = options
	if cfgOpts.Viper != nil {
		viper = cfgOpts.Viper
	} else {
		viper = globalViper.GetViper()
	}
	cfg = retag.ConvertAny(config, mapstructureRetagger{})

	if rootCmd == nil { // only viper
		return initializeViper()
	}
	cobra.OnInitialize(func() {
		err := initializeViper()
		if err != nil {
			panic(fmt.Errorf("panic from WONSZ lib: cannot initialize viper: %w", err))
		}
	})

	confType := reflect.TypeOf(cfg).Elem()
	return bindFieldsRecursive(rootCmd.PersistentFlags(), confType, "", "")
}

func bindFieldsRecursive(flags *pflag.FlagSet, t reflect.Type, namePrefix, mappingPrefix string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}

		if field.Tag.Get("wonsz-flag-ignore") != "" {
			continue
		}

		dashedName := camelCaseToDashedLowered(field.Name)
		underscoredName := camelCaseToUnderscoredLowered(field.Name)
		mappingName := underscoredName

		if namePrefix != "" {
			dashedName = namePrefix + "-" + dashedName
		}
		if mappingPrefix != "" {
			mappingName = mappingPrefix + "." + underscoredName
		}

		// Handle nested structs (excluding special types like time.Time)
		if field.Type.Kind() == reflect.Struct &&
			field.Type != reflect.TypeOf(time.Time{}) &&
			field.Type != reflect.TypeOf(net.IP{}) &&
			field.Type != reflect.TypeOf(net.IPNet{}) {
			if err := bindFieldsRecursive(flags, field.Type, dashedName, mappingName); err != nil {
				return err
			}
			continue
		}

		usageHint := field.Tag.Get("usage")
		shortcut, _ := field.Tag.Lookup("shortcut")
		err := bindPFlag(flags, field, dashedName, shortcut, usageHint)
		if err != nil {
			if cfgOpts.IgnoreViperBindErrors {
				continue
			}
			return fmt.Errorf("cannot bind flag %s: %w. "+
				"You can ignore this error by setting IgnoreViperBindErrors to true"+
				"or adding wonsz-flag-ignore annotation to field", dashedName, err)
		}

		targetFlag := flags.Lookup(dashedName)
		if targetFlag == nil {
			return fmt.Errorf("flag %s not found, despite successful binding", dashedName)
		}
		err = viper.BindPFlag(mappingName, targetFlag)
		if err != nil {
			return err
		}
	}
	return nil
}

func initializeViper() error {
	viper.SetEnvPrefix(cfgOpts.EnvPrefix)

	for _, path := range cfgOpts.ConfigPaths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigType(cfgOpts.ConfigType)
	viper.SetConfigName(cfgOpts.ConfigName)

	err := bindEnvsAndSetDefaults()
	if err != nil {
		return fmt.Errorf("cannot bind env variables, err: %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if cfgOpts.ConfigName != "" {
			return fmt.Errorf("cannot read config file: %w", err)
		}
	}

	if err = viper.Unmarshal(&cfg, globalViper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		mapstructure.StringToTimeHookFunc(time.RFC3339),
	))); err != nil {
		return fmt.Errorf("cannot unmarshal config into config struct: %w", err)
	}
	return nil
}

func bindEnvsAndSetDefaults() error {
	el := reflect.TypeOf(cfg).Elem()
	return processStructFields(el, "")
}

func processStructFields(t reflect.Type, prefix string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		mapping := field.Tag.Get("mapstructure")
		if field.Anonymous || mapping == "" {
			continue
		}
		if prefix != "" {
			mapping = prefix + "." + mapping
		}

		if field.Type.Kind() == reflect.Struct {
			if field.Type.String() != "time.Time" {
				err := processStructFields(field.Type, mapping)
				if err != nil {
					return err
				}
				continue
			}
		}

		defaultVal := field.Tag.Get("default")
		if defaultVal != "" {
			viper.SetDefault(mapping, defaultVal)
		} else {
			err := viper.BindEnv(mapping)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func bindPFlag(flags *pflag.FlagSet, field reflect.StructField, dashedName, shortcut, usageHint string) error {
	switch field.Type.Kind() {
	case reflect.String:
		flags.StringP(dashedName, shortcut, "", usageHint)
	case reflect.Int64:
		if field.Type == reflect.TypeOf(time.Duration(0)) {
			flags.DurationP(dashedName, shortcut, 0, usageHint)
		} else {
			flags.Int64P(dashedName, shortcut, 0, usageHint)
		}
	case reflect.Int32:
		flags.Int32P(dashedName, shortcut, 0, usageHint)
	case reflect.Int16:
		flags.Int16P(dashedName, shortcut, 0, usageHint)
	case reflect.Int8:
		flags.Int8P(dashedName, shortcut, 0, usageHint)
	case reflect.Int:
		flags.IntP(dashedName, shortcut, 0, usageHint)
	case reflect.Uint:
		flags.UintP(dashedName, shortcut, 0, usageHint)
	case reflect.Uint64:
		flags.Uint64P(dashedName, shortcut, 0, usageHint)
	case reflect.Uint32:
		flags.Uint32P(dashedName, shortcut, 0, usageHint)
	case reflect.Uint16:
		flags.Uint16P(dashedName, shortcut, 0, usageHint)
	case reflect.Uint8:
		flags.Uint8P(dashedName, shortcut, 0, usageHint)
	case reflect.Float64:
		flags.Float64P(dashedName, shortcut, 0, usageHint)
	case reflect.Float32:
		flags.Float32P(dashedName, shortcut, 0, usageHint)
	case reflect.Bool:
		flags.BoolP(dashedName, shortcut, false, usageHint)
	case reflect.Struct:
		switch field.Type {
		case reflect.TypeOf(time.Time{}):
			flags.TimeP(dashedName, shortcut, time.Time{}, []string{time.RFC3339}, usageHint)
		case reflect.TypeOf(net.IP{}):
			flags.IPP(dashedName, shortcut, net.IP{}, usageHint)
		case reflect.TypeOf(net.IPNet{}):
			flags.IPNetP(dashedName, shortcut, net.IPNet{}, usageHint)
		default:
			return fmt.Errorf("unsupported flag %s type: %s", dashedName, field.Type.String())
		}
	case reflect.Array:
		if field.Type.Elem().Kind() == reflect.String {
			flags.StringArrayP(dashedName, shortcut, []string{}, usageHint)
		} else {
			return fmt.Errorf("unsupported flag %s type: %s. only string arrays are supported",
				dashedName, field.Type.String())
		}
	case reflect.Slice:
		switch field.Type.Elem().Kind() {
		case reflect.String:
			flags.StringSliceP(dashedName, shortcut, []string{}, usageHint)
		case reflect.Int:
			flags.IntSliceP(dashedName, shortcut, []int{}, usageHint)
		case reflect.Int32:
			flags.Int32SliceP(dashedName, shortcut, []int32{}, usageHint)
		case reflect.Int64:
			flags.Int64SliceP(dashedName, shortcut, []int64{}, usageHint)
		case reflect.Uint:
			flags.UintSliceP(dashedName, shortcut, []uint{}, usageHint)
		case reflect.Uint8:
			flags.BytesHexP(dashedName, shortcut, []byte{}, usageHint)
		case reflect.Float32:
			flags.Float32SliceP(dashedName, shortcut, []float32{}, usageHint)
		case reflect.Float64:
			flags.Float64SliceP(dashedName, shortcut, []float64{}, usageHint)
		default:
			return fmt.Errorf("unsupported slice flag %s type: %s",
				dashedName, field.Type.String())
		}
	case reflect.Map:
		switch field.Type.Elem().Kind() {
		case reflect.String:
			flags.StringToStringP(dashedName, shortcut, map[string]string{}, usageHint)
		case reflect.Int:
			flags.StringToIntP(dashedName, shortcut, map[string]int{}, usageHint)
		case reflect.Int32:
			flags.StringToIntP(dashedName, shortcut, map[string]int{}, usageHint)
		case reflect.Int64:
			flags.StringToInt64P(dashedName, shortcut, map[string]int64{}, usageHint)
		default:
			return fmt.Errorf("unsupported flag %s type: %s",
				dashedName, field.Type.String())
		}
	default:
		return fmt.Errorf("unsupported flag %s type: %s", dashedName, field.Type.String())
	}
	return nil
}
