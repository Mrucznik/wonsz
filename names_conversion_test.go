package wonsz

import (
	"reflect"
	"testing"
)

func Test_camelCaseToDashedLowered(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelCaseToDashedLowered(tt.args.text); got != tt.want {
				t.Errorf("camelCaseToDashedLowered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_camelCaseToSeparatorsLowered(t *testing.T) {
	type args struct {
		text      string
		separator rune
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelCaseToSeparatorsLowered(tt.args.text, tt.args.separator); got != tt.want {
				t.Errorf("camelCaseToSeparatorsLowered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_camelCaseToUnderscoredLowered(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelCaseToUnderscoredLowered(tt.args.text); got != tt.want {
				t.Errorf("camelCaseToUnderscoredLowered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDesiredSeparatorPositions(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want map[int]struct{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDesiredSeparatorPositions(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDesiredSeparatorPositions() = %v, want %v", got, tt.want)
			}
		})
	}
}
