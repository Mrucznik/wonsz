package wonsz

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/sevlyar/retag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	globalViper "github.com/spf13/viper"
)

// Wonsz binds a single configuration struct to a config file, environment
// variables and cobra command flags. Create instances with New; each instance
// keeps its own options and viper, so multiple configs can coexist.
type Wonsz struct {
	opts        ConfigOpts
	cfg         interface{} // retagged copy sharing memory with originalCfg
	originalCfg interface{}
	viper       *globalViper.Viper
}

// defaultWonsz backs the package-level BindConfig/Get/GetViper API.
var defaultWonsz *Wonsz

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

	// If true, Wonsz will not return an error if a config field cannot be bound to a flag
	// or the resulting flag cannot be bound to viper. Such fields are silently skipped.
	IgnoreViperBindErrors bool

	// If true, Wonsz watches the config file and re-unmarshals the config struct
	// when the file changes. Note that the struct is updated from a background
	// goroutine, so guard access to it if your application reads it concurrently.
	WatchConfig bool
}

// Get returns the config struct instance passed to BindConfig.
// The result can be type-asserted back to the original pointer type.
func Get() interface{} {
	if defaultWonsz == nil {
		return nil
	}
	return defaultWonsz.Get()
}

// GetViper returns a viper instance used by Wonsz.
func GetViper() *globalViper.Viper {
	if defaultWonsz == nil {
		return nil
	}
	return defaultWonsz.Viper()
}

// BindConfig binds configuration structure to config file, environment variables and cobra command flags.
// The config parameter should be a pointer to the configuration structure.
// You can pass nil to rootCmd if you don't want to bind cobra command flags with config.
// It uses a package-level default instance; to bind multiple independent configs use New.
func BindConfig(config interface{}, rootCmd *cobra.Command, options ConfigOpts) error {
	w, err := New(config, rootCmd, options)
	if w != nil {
		defaultWonsz = w
	}
	return err
}

// New binds the configuration structure like BindConfig, but returns an
// independent Wonsz instance instead of using package-level state.
func New(config interface{}, rootCmd *cobra.Command, options ConfigOpts) (*Wonsz, error) {
	if reflect.TypeOf(config).Kind() != reflect.Ptr || reflect.TypeOf(config).Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("config parameter is not a pointer to a structure. Maybe you should use & operator")
	}

	// prepare for processing
	w := &Wonsz{opts: options, originalCfg: config}
	if options.Viper != nil {
		w.viper = options.Viper
	} else {
		w.viper = globalViper.GetViper()
	}
	w.cfg = retag.ConvertAny(config, mapstructureRetagger{})

	if rootCmd == nil { // only viper
		return w, w.initializeViper()
	}
	cobra.OnInitialize(func() {
		err := w.initializeViper()
		if err != nil {
			panic(fmt.Errorf("panic from WONSZ lib: cannot initialize viper: %w", err))
		}
	})

	confType := reflect.TypeOf(w.cfg).Elem()
	return w, w.bindFieldsRecursive(rootCmd.PersistentFlags(), confType, "", "")
}

// Get returns the config struct instance passed to New.
// The result can be type-asserted back to the original pointer type.
func (w *Wonsz) Get() interface{} {
	return w.originalCfg
}

// Viper returns the viper instance used by this Wonsz instance.
func (w *Wonsz) Viper() *globalViper.Viper {
	return w.viper
}

func (w *Wonsz) bindFieldsRecursive(flags *pflag.FlagSet, t reflect.Type, namePrefix, mappingPrefix string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}

		if field.Tag.Get("wonsz-flag-ignore") == "true" {
			continue
		}

		// The retagged config type always carries a mapstructure tag; a custom tag
		// set by the user takes precedence over the name derived from the field name.
		underscoredName := strings.Split(field.Tag.Get("mapstructure"), ",")[0]
		if underscoredName == "-" {
			continue
		}
		if underscoredName == "" {
			underscoredName = camelCaseToUnderscoredLowered(field.Name)
		}
		dashedName := strings.ReplaceAll(underscoredName, "_", "-")
		mappingName := underscoredName

		if namePrefix != "" {
			dashedName = namePrefix + "-" + dashedName
		}
		if mappingPrefix != "" {
			mappingName = mappingPrefix + "." + underscoredName
		}

		// Handle nested structs, also behind pointers (excluding special types like time.Time)
		nestedType := field.Type
		if nestedType.Kind() == reflect.Ptr {
			nestedType = nestedType.Elem()
		}
		if isNestedStruct(nestedType) {
			if err := w.bindFieldsRecursive(flags, nestedType, dashedName, mappingName); err != nil {
				return err
			}
			continue
		}

		usageHint := field.Tag.Get("usage")
		shortcut, _ := field.Tag.Lookup("shortcut")
		if len(shortcut) > 1 {
			return fmt.Errorf("invalid shortcut %q for field %s: must be a single ASCII character",
				shortcut, field.Name)
		}
		err := bindPFlag(flags, field, dashedName, shortcut, usageHint)
		if err != nil {
			if w.opts.IgnoreViperBindErrors {
				continue
			}
			return fmt.Errorf("cannot bind flag %s: %w. "+
				"You can ignore this error by setting IgnoreViperBindErrors to true "+
				"or by adding the wonsz-flag-ignore annotation to the field", dashedName, err)
		}

		targetFlag := flags.Lookup(dashedName)
		if targetFlag == nil {
			return fmt.Errorf("flag %s not found, despite successful binding", dashedName)
		}
		err = w.viper.BindPFlag(mappingName, targetFlag)
		if err != nil {
			if w.opts.IgnoreViperBindErrors {
				continue
			}
			return fmt.Errorf("cannot bind flag %s to viper key %s: %w", dashedName, mappingName, err)
		}
	}
	return nil
}

// isNestedStruct reports whether t is a struct that should be recursed into
// when binding fields, as opposed to leaf types like time.Time or net.IPNet
// that are bound as single values.
func isNestedStruct(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	// Compare with ConvertibleTo instead of type identity, because the retagged
	// config type may contain structurally identical copies of these types
	// (struct tags are ignored in convertibility checks).
	if t.ConvertibleTo(reflect.TypeOf(time.Time{})) || t.ConvertibleTo(reflect.TypeOf(net.IPNet{})) {
		return false
	}
	return true
}

// stringToRetaggedIPNetHookFunc works like mapstructure.StringToIPNetHookFunc,
// but also matches the structurally identical copy of net.IPNet that retag
// produces inside the retagged config type.
func stringToRetaggedIPNetHookFunc() mapstructure.DecodeHookFuncType {
	ipNetType := reflect.TypeOf(net.IPNet{})
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Struct || !t.ConvertibleTo(ipNetType) {
			return data, nil
		}
		_, ipNet, err := net.ParseCIDR(data.(string))
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(*ipNet).Convert(t).Interface(), nil
	}
}

func (w *Wonsz) initializeViper() error {
	w.viper.SetEnvPrefix(w.opts.EnvPrefix)

	for _, path := range w.opts.ConfigPaths {
		w.viper.AddConfigPath(path)
	}
	w.viper.SetConfigType(w.opts.ConfigType)
	w.viper.SetConfigName(w.opts.ConfigName)

	err := w.bindEnvsAndSetDefaults()
	if err != nil {
		return fmt.Errorf("cannot bind env variables, err: %w", err)
	}

	w.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	w.viper.AutomaticEnv()

	if err = w.viper.ReadInConfig(); err != nil {
		if w.opts.ConfigName != "" {
			return fmt.Errorf("cannot read config file: %w", err)
		}
	}

	if w.opts.WatchConfig {
		w.viper.OnConfigChange(func(fsnotify.Event) {
			_ = w.unmarshalConfig()
		})
		w.viper.WatchConfig()
	}

	return w.unmarshalConfig()
}

func (w *Wonsz) unmarshalConfig() error {
	if err := w.viper.Unmarshal(&w.cfg, globalViper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToIPHookFunc(),
		stringToRetaggedIPNetHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		mapstructure.StringToTimeHookFunc(time.RFC3339),
	))); err != nil {
		return fmt.Errorf("cannot unmarshal config into config struct: %w", err)
	}
	return nil
}

func (w *Wonsz) bindEnvsAndSetDefaults() error {
	el := reflect.TypeOf(w.cfg).Elem()
	return w.processStructFields(el, "")
}

func (w *Wonsz) processStructFields(t reflect.Type, prefix string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		mapping := strings.Split(field.Tag.Get("mapstructure"), ",")[0]
		if field.Anonymous || mapping == "" || mapping == "-" {
			continue
		}
		if prefix != "" {
			mapping = prefix + "." + mapping
		}

		nestedType := field.Type
		if nestedType.Kind() == reflect.Ptr {
			nestedType = nestedType.Elem()
		}
		if isNestedStruct(nestedType) {
			err := w.processStructFields(nestedType, mapping)
			if err != nil {
				return err
			}
			continue
		}

		defaultVal := field.Tag.Get("default")
		if defaultVal != "" {
			w.viper.SetDefault(mapping, defaultVal)
		} else {
			err := w.viper.BindEnv(mapping)
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
		switch {
		case field.Type.ConvertibleTo(reflect.TypeOf(time.Time{})):
			flags.TimeP(dashedName, shortcut, time.Time{}, []string{time.RFC3339}, usageHint)
		case field.Type.ConvertibleTo(reflect.TypeOf(net.IPNet{})):
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
			if field.Type == reflect.TypeOf(net.IP{}) {
				flags.IPP(dashedName, shortcut, net.IP{}, usageHint)
			} else {
				flags.BytesHexP(dashedName, shortcut, []byte{}, usageHint)
			}
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
