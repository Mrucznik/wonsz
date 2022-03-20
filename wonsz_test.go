package wonsz

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"testing"
)

func TestInitializeConfig(t *testing.T) {
	var testConfig struct {
		AlwaysInConfig string
		TestField      string
	}

	err := BindConfig(&testConfig, &cobra.Command{}, ConfigOpts{})
	if err != nil {
		t.Fatal(err)
	}
}

func ExampleBindConfig() {
	os.Setenv("EXAMPLE_FIELD", "this is my example config field")

	var myConfig struct {
		ExampleField string
	}

	err := BindConfig(&myConfig, nil, ConfigOpts{})
	if err != nil {
		panic(err)
	}

	fmt.Println(myConfig.ExampleField)
	// Output: this is my example config field
}

func Test_BindConfig_withEnv(t *testing.T) {
	os.Setenv("SLICE_FIELD", "some,text,here")

	var testConfig struct {
		SliceField []string
	}

	err := BindConfig(&testConfig, nil, ConfigOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if testConfig.SliceField[0] != "some" ||
		testConfig.SliceField[1] != "text" ||
		testConfig.SliceField[2] != "here" {
		t.Errorf("Expected %s, got %s", "some", testConfig.SliceField[0])
	}
}

func Test_BindConfig_withFlag(t *testing.T) {
	var testConfig struct {
		SliceField []string
	}

	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			if testConfig.SliceField[0] != "some" ||
				testConfig.SliceField[1] != "text" ||
				testConfig.SliceField[2] != "here" {
				t.Errorf("Expected %s, got %s", "some", testConfig.SliceField[0])
			}
		},
	}

	os.Args = []string{"cmd", "--slice-field=some,text,here"}

	err := BindConfig(&testConfig, rootCmd, ConfigOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
}
