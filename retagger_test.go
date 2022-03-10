package wonsz

import (
	"reflect"
	"testing"
)

func Test_mapstructureRetagger_MakeTag(t *testing.T) {
	type testStruct struct {
		testField string
	}

	type args struct {
		structureType reflect.Type
		fieldIndex    int
	}
	tests := []struct {
		name string
		args args
		want reflect.StructTag
	}{
		{
			name: "basic test",
			args: args{
				structureType: reflect.TypeOf(testStruct{}),
				fieldIndex:    0,
			},
			want: "mapstructure:\"test_field\"",
		},
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
