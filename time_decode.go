package wonsz

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(time.Duration(0)) {
			switch f.Kind() {
			case reflect.String:
				return time.ParseDuration(data.(string))
			case reflect.Float64:
				return time.Duration(data.(float64)) * time.Millisecond, nil
			case reflect.Int64:
				return time.Duration(data.(int64)) * time.Millisecond, nil
			default:
				return data, nil
			}
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		timeFormats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.RFC1123,
			time.RFC1123Z,
		}

		switch f.Kind() {
		case reflect.String:
			str := data.(string)
			for _, layout := range timeFormats {
				if t, err := time.Parse(layout, str); err == nil {
					return t, nil
				}
			}
			return time.Time{}, fmt.Errorf("unable to parse time string: %s", str)
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

func MapstructureDecoder() func(*mapstructure.DecoderConfig) {
	return func(config *mapstructure.DecoderConfig) {
		config.DecodeHook = ToTimeHookFunc()
	}
}
