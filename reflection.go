package wonsz

import (
	"reflect"
	"strconv"
)

// Tag contains the name and value of a structure field tag.
type Tag struct {
	name  string
	value string
}

// GetTagsForField extracts field tag names and values from specified structure.
func GetTagsForField(field reflect.StructField) []Tag {
	tag := field.Tag
	var tags []Tag

	// this block of code is modification of https://github.com/golang/go/blob/master/src/reflect/type.go#L1191-L1238
	// Copyright (c) 2009 The Go Authors. https://github.com/golang/go/blob/master/LICENSE
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		tags = append(tags, Tag{name: name, value: value})
	}
	return tags
}
