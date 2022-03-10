package wonsz

import (
	"reflect"
	"testing"
)

func TestGetTagsForField(t *testing.T) {
	type testStruct struct {
		noTags   string
		oneTag   string `someTag:"testTagValue"`
		manyTags string `tag1:"value1" tag2:"value2" tag3:"value3"`
	}

	type args struct {
		field reflect.StructField
	}
	tests := []struct {
		name string
		args args
		want []Tag
	}{
		{
			name: "no tags",
			args: args{
				field: reflect.TypeOf(testStruct{}).Field(0),
			},
			want: nil,
		},
		{
			name: "one tag",
			args: args{
				field: reflect.TypeOf(testStruct{}).Field(1),
			},
			want: []Tag{{name: "someTag", value: "testTagValue"}},
		},
		{
			name: "many tags",
			args: args{
				field: reflect.TypeOf(testStruct{}).Field(2),
			},
			want: []Tag{
				{name: "tag1", value: "value1"},
				{name: "tag2", value: "value2"},
				{name: "tag3", value: "value3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTagsForField(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsForField() = %v, want %v", got, tt.want)
			}
		})
	}
}
