package wonsz

import (
	"fmt"
	"reflect"
)

// Structure for being used by retag.Convert to create desired mapstructure tags for config structure.
// It's needed by viper to be able to map environment variables.
type mapstructureRetagger struct{}

// Add a mapstructure tag to structure field.
func (m mapstructureRetagger) MakeTag(structureType reflect.Type, fieldIndex int) reflect.StructTag {
	field := structureType.Field(fieldIndex)
	mapping := camelCaseToUnderscoredLowered(field.Name)

	tags := GetTagsForField(field)
	for i := range tags {
		if tags[i].name == "mapstructure" {
			return field.Tag
		}
	}

	var newTag string
	if field.Tag != "" {
		newTag = fmt.Sprintf("mapstructure:\"%s\" %s", mapping, field.Tag)
	} else {
		newTag = fmt.Sprintf("mapstructure:\"%s\"", mapping)
	}
	return reflect.StructTag(newTag)
}
