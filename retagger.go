package main

import (
	"fmt"
	"reflect"
	"strings"
)

type MapstructureRetagger struct{}

func (m MapstructureRetagger) MakeTag(structureType reflect.Type, fieldIndex int) reflect.StructTag {
	field := structureType.Field(fieldIndex)
	mapping := CamelCaseToUnderscoredLowered(field.Name)
	newTag := strings.Join([]string{fmt.Sprintf("mapstructure:\"%s\"", mapping), string(field.Tag)}, " ")
	return reflect.StructTag(newTag)
}
