package wonsz

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	globalViper "github.com/spf13/viper"
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
	_ = os.Setenv("EXAMPLE_FIELD", "this is my example config field")

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
	_ = os.Setenv("SLICE_FIELD", "some,text,here")

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

func Test_BindConfig_unsupportedTypeErrorMessage(t *testing.T) {
	var testConfig struct {
		Unsupported complex128
	}

	err := BindConfig(&testConfig, &cobra.Command{}, ConfigOpts{})
	if err == nil {
		t.Fatal("expected an error for unsupported field type, got nil")
	}
	if !strings.Contains(err.Error(), "IgnoreViperBindErrors to true or by adding") {
		t.Errorf("malformed error message: %q", err.Error())
	}
}

func Test_BindConfig_flagIgnoreTagValue(t *testing.T) {
	var testConfig struct {
		Kept    string `wonsz-flag-ignore:"false"`
		Skipped string `wonsz-flag-ignore:"true"`
	}

	cmd := &cobra.Command{}
	err := BindConfig(&testConfig, cmd, ConfigOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if cmd.PersistentFlags().Lookup("kept") == nil {
		t.Error(`field with wonsz-flag-ignore:"false" should be bound to a flag`)
	}
	if cmd.PersistentFlags().Lookup("skipped") != nil {
		t.Error(`field with wonsz-flag-ignore:"true" should not be bound to a flag`)
	}
}

func Test_BindConfig_customMapstructureTag(t *testing.T) {
	var testConfig struct {
		FieldOne string `mapstructure:"custom_name"`
	}

	cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	err := BindConfig(&testConfig, cmd, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	if cmd.PersistentFlags().Lookup("custom-name") == nil {
		t.Fatal("flag name should be derived from the mapstructure tag")
	}

	cmd.SetArgs([]string{"--custom-name", "from-flag"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if testConfig.FieldOne != "from-flag" {
		t.Errorf("FieldOne: got %q, want %q", testConfig.FieldOne, "from-flag")
	}
}

func Test_BindConfig_ipFieldsFromEnv(t *testing.T) {
	t.Setenv("SERVER_IP", "192.168.1.10")
	t.Setenv("SERVER_SUBNET", "10.0.0.0/8")

	var testConfig struct {
		ServerIp     net.IP
		ServerSubnet net.IPNet
	}

	err := BindConfig(&testConfig, nil, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := testConfig.ServerIp.String(), "192.168.1.10"; got != want {
		t.Errorf("ServerIp: got %q, want %q", got, want)
	}
	if got, want := testConfig.ServerSubnet.String(), "10.0.0.0/8"; got != want {
		t.Errorf("ServerSubnet: got %q, want %q", got, want)
	}
}

func Test_BindConfig_ipFieldsFromFlags(t *testing.T) {
	var testConfig struct {
		ServerIp     net.IP
		ServerSubnet net.IPNet
	}

	cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	err := BindConfig(&testConfig, cmd, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	cmd.SetArgs([]string{"--server-ip", "172.16.0.1", "--server-subnet", "192.168.0.0/16"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if got, want := testConfig.ServerIp.String(), "172.16.0.1"; got != want {
		t.Errorf("ServerIp: got %q, want %q", got, want)
	}
	if got, want := testConfig.ServerSubnet.String(), "192.168.0.0/16"; got != want {
		t.Errorf("ServerSubnet: got %q, want %q", got, want)
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
