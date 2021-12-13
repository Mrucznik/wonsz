package wonsz

import (
	"reflect"
	"testing"
)

func Test_mapstructureRetagger_MakeTag(t *testing.T) {
	type args struct {
		structureType reflect.Type
		fieldIndex    int
	}
	tests := []struct {
		name string
		args args
		want reflect.StructTag
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mapstructureRetagger{}
			if got := m.MakeTag(tt.args.structureType, tt.args.fieldIndex); got != tt.want {
				t.Errorf("MakeTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
