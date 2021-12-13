package wonsz

import (
	"reflect"
	"testing"
)

func TestGetTagsForField(t *testing.T) {
	type args struct {
		field reflect.StructField
	}
	tests := []struct {
		name string
		args args
		want []Tag
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTagsForField(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsForField() = %v, want %v", got, tt.want)
			}
		})
	}
}
