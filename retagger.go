package wonsz

import (
	"fmt"
	"reflect"
	"strings"
)

// Here some description
type mapstructureRetagger struct{}

func (m mapstructureRetagger) MakeTag(structureType reflect.Type, fieldIndex int) reflect.StructTag {
	field := structureType.Field(fieldIndex)
	mapping := camelCaseToUnderscoredLowered(field.Name)
	newTag := strings.Join([]string{fmt.Sprintf("mapstructure:\"%s\"", mapping), string(field.Tag)}, " ")
	return reflect.StructTag(newTag)
}
