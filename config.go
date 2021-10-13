package wonsz

import (
	"fmt"
	"reflect"
)

type IConfig interface {
}

type Config struct {
	AlwaysInConfig string `sometag:"xd" boolTag:""`
}

type ConfigOpts struct {
}

func InitializeConfig(config IConfig, _ *ConfigOpts) {
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
		//fmt.Printf("field %d: %s\n", i, field.Name)

		if field.Anonymous {
			continue
		}

		// get all tags
		for j, key := range GetTagsForField(field) {
			fmt.Printf("tag %d: %+v\n", j, key)
		}

		dashed, underscored := ConvertFromCamelCase(field.Name)

		fmt.Printf("----\noriginal: %s\nunderscored: %s\ndashed: %s\n", field.Name, string(underscored), string(dashed))
	}
}
