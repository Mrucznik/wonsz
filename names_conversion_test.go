package wonsz

import (
	"fmt"
	"testing"
)

func Example_getDesiredSeparatorPositions() {
	separatorPos := getDesiredSeparatorPositions("userIDText")
	fmt.Println(separatorPos)

	// Output: map[4:{} 6:{}]
}

func Test_camelCaseToDashedLowered(t *testing.T) {
	tests := []struct {
		name          string
		textToConvert string
		want          string
	}{
		{
			name:          "simpleTestingName",
			textToConvert: "simpleTestingName",
			want:          "simple-testing-name",
		},
		{
			name:          "longCAPITALWord",
			textToConvert: "longCAPITALWord",
			want:          "long-capital-word",
		},
		{
			name:          "non letters symbols",
			textToConvert: "Non-ConventionalStringMyID",
			want:          "non-conventional-string-my-id",
		},
		{
			name:          "VeIrDTeXt",
			textToConvert: "VeIrDTeXt",
			want:          "ve-ir-d-te-xt",
		},
		{
			name:          "empty",
			textToConvert: "",
			want:          "",
		},
		{
			name:          "numbers in name",
			textToConvert: "userWith5BANSButNoMoreTHAN100Ok1200ok",
			want:          "user-with-5-bans-but-no-more-than-100-ok-1200-ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelCaseToDashedLowered(tt.textToConvert); got != tt.want {
				t.Errorf("camelCaseToDashedLowered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_camelCaseToUnderscoredLowered(t *testing.T) {
	tests := []struct {
		name          string
		textToConvert string
		want          string
	}{

		{
			name:          "simpleTestingName",
			textToConvert: "simpleTestingName",
			want:          "simple_testing_name",
		},
		{
			name:          "longCAPITALWord",
			textToConvert: "longCAPITALWord",
			want:          "long_capital_word",
		},
		{
			name:          "non letters symbols",
			textToConvert: "Non-ConventionalStringMyID",
			want:          "non-conventional_string_my_id",
		},
		{
			name:          "VeIrDTeXt",
			textToConvert: "VeIrDTeXt",
			want:          "ve_ir_d_te_xt",
		},
		{
			name:          "empty",
			textToConvert: "",
			want:          "",
		},
		{
			name:          "numbers in name",
			textToConvert: "userWith5BANSButNoMoreTHAN100Ok1200ok",
			want:          "user_with_5_bans_but_no_more_than_100_ok_1200_ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelCaseToUnderscoredLowered(tt.textToConvert); got != tt.want {
				t.Errorf("camelCaseToUnderscoredLowered() = %v, want %v", got, tt.want)
			}
		})
	}
}
