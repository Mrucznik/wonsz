package wonsz

import (
	"fmt"
	"reflect"
	"strings"
)

// Structure for being used by retag.Convert to create desired mapstructure tags for config structure.
type mapstructureRetagger struct{}

// Add a mapstructure tag to structure field.
func (m mapstructureRetagger) MakeTag(structureType reflect.Type, fieldIndex int) reflect.StructTag {
	field := structureType.Field(fieldIndex)
	mapping := camelCaseToUnderscoredLowered(field.Name)
	newTag := strings.Join([]string{fmt.Sprintf("mapstructure:\"%s\"", mapping), string(field.Tag)}, " ")
	return reflect.StructTag(newTag)
}
