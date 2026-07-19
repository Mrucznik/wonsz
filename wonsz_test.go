package wonsz

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	if !strings.Contains(err.Error(), "IgnoreFlagBindErrors to true or by adding") {
		t.Errorf("malformed error message: %q", err.Error())
	}
}

func Test_BindConfig_flagIgnoreTag(t *testing.T) {
	t.Setenv("SKIPPED", "from-env")

	var testConfig struct {
		Kept    string
		Skipped string `wonsz:"flag-ignore"`
	}

	cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	err := BindConfig(&testConfig, cmd, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	if cmd.PersistentFlags().Lookup("kept") == nil {
		t.Error("field without ignore tag should be bound to a flag")
	}
	if cmd.PersistentFlags().Lookup("skipped") != nil {
		t.Error(`field with wonsz:"flag-ignore" should not be bound to a flag`)
	}

	// env binding must still work for a flag-ignored field
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if testConfig.Skipped != "from-env" {
		t.Errorf("Skipped: got %q, want %q (env should still bind)", testConfig.Skipped, "from-env")
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

func TestGetReturnsTypedPointer(t *testing.T) {
	t.Setenv("MY_FIELD", "from-env")

	type conf struct{ MyField string }
	var c conf

	w, err := New(&c, nil, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	got := w.Get() // typed *conf, no type assertion needed
	if got != &c {
		t.Error("Get() should return the same pointer that was passed to New")
	}
	if got.MyField != "from-env" {
		t.Errorf("MyField: got %q, want %q", got.MyField, "from-env")
	}
}

func TestNewRejectsNonStruct(t *testing.T) {
	var notAStruct int
	if _, err := New(&notAStruct, nil, ConfigOpts{Viper: globalViper.New()}); err == nil {
		t.Error("expected an error for a non-struct config type, got nil")
	}

	var nilConfig *struct{ Field string }
	if _, err := New(nilConfig, nil, ConfigOpts{Viper: globalViper.New()}); err == nil {
		t.Error("expected an error for a nil config pointer, got nil")
	}
}

func Test_BindConfig_invalidShortcut(t *testing.T) {
	var testConfig struct {
		Field string `shortcut:"ab"`
	}

	err := BindConfig(&testConfig, &cobra.Command{}, ConfigOpts{Viper: globalViper.New()})
	if err == nil {
		t.Fatal("expected an error for a multi-character shortcut, got nil")
	}
	if !strings.Contains(err.Error(), "shortcut") {
		t.Errorf("error should mention the shortcut tag, got: %q", err.Error())
	}
}

func Test_BindConfig_pointerToStructField(t *testing.T) {
	type Inner struct {
		Port int
	}
	var testConfig struct {
		Server *Inner
	}

	cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	err := BindConfig(&testConfig, cmd, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	cmd.SetArgs([]string{"--server-port", "9090"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if testConfig.Server == nil {
		t.Fatal("Server: got nil, want allocated struct")
	}
	if testConfig.Server.Port != 9090 {
		t.Errorf("Server.Port: got %d, want %d", testConfig.Server.Port, 9090)
	}
}

func Test_BindConfig_defaultTagWithEnv(t *testing.T) {
	type conf struct {
		Field string `default:"default-value"`
	}

	var withDefault conf
	if err := BindConfig(&withDefault, nil, ConfigOpts{Viper: globalViper.New()}); err != nil {
		t.Fatal(err)
	}
	if withDefault.Field != "default-value" {
		t.Errorf("Field: got %q, want %q", withDefault.Field, "default-value")
	}

	t.Setenv("FIELD", "env-value")
	var withEnv conf
	if err := BindConfig(&withEnv, nil, ConfigOpts{Viper: globalViper.New()}); err != nil {
		t.Fatal(err)
	}
	if withEnv.Field != "env-value" {
		t.Errorf("Field: got %q, want %q (env should override the default tag)", withEnv.Field, "env-value")
	}
}

func Test_BindConfig_wonszDashExcludesField(t *testing.T) {
	t.Setenv("EXCLUDED", "from-env")
	t.Setenv("INCLUDED", "from-env")

	var testConfig struct {
		Excluded string `wonsz:"-"`
		Included string
	}

	cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	err := BindConfig(&testConfig, cmd, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	if cmd.PersistentFlags().Lookup("excluded") != nil {
		t.Error(`field with wonsz:"-" should not be bound to a flag`)
	}

	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if testConfig.Excluded != "" {
		t.Errorf(`Excluded: got %q, want "" (field with wonsz:"-" should not be populated)`, testConfig.Excluded)
	}
	if testConfig.Included != "from-env" {
		t.Errorf("Included: got %q, want %q", testConfig.Included, "from-env")
	}
}

func Test_BindConfig_watchConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(`{"watched_field": "initial"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var testConfig struct {
		WatchedField string
	}

	// The reload happens on a background goroutine, so the config value is read
	// inside the OnConfigChange callback (same goroutine as the re-unmarshal)
	// and handed to the test through a channel.
	reloaded := make(chan string, 8)

	err := BindConfig(&testConfig, nil, ConfigOpts{
		ConfigPaths: []string{dir},
		ConfigType:  "json",
		ConfigName:  "config",
		Viper:       globalViper.New(),
		WatchConfig: true,
		OnConfigChange: func(err error) {
			if err == nil {
				reloaded <- testConfig.WatchedField
			}
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if testConfig.WatchedField != "initial" {
		t.Fatalf("WatchedField: got %q, want %q", testConfig.WatchedField, "initial")
	}

	if err := os.WriteFile(path, []byte(`{"watched_field": "updated"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	deadline := time.After(3 * time.Second)
	for {
		select {
		case got := <-reloaded:
			if got == "updated" {
				return
			}
		case <-deadline:
			t.Fatal("timed out waiting for config reload after file change")
		}
	}
}

func TestNewIndependentInstances(t *testing.T) {
	t.Setenv("FIRST_FIELD", "first")
	t.Setenv("SECOND_FIELD", "second")

	type confA struct{ FirstField string }
	type confB struct{ SecondField string }
	var a confA
	var b confB

	cmdA := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	cmdB := &cobra.Command{Run: func(*cobra.Command, []string) {}}

	wA, err := New(&a, cmdA, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}
	wB, err := New(&b, cmdB, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	// Deferred (cobra.OnInitialize) initialization of the first instance must
	// not be clobbered by the second New call.
	cmdA.SetArgs([]string{})
	if err := cmdA.Execute(); err != nil {
		t.Fatal(err)
	}
	cmdB.SetArgs([]string{})
	if err := cmdB.Execute(); err != nil {
		t.Fatal(err)
	}

	if a.FirstField != "first" {
		t.Errorf("FirstField: got %q, want %q", a.FirstField, "first")
	}
	if b.SecondField != "second" {
		t.Errorf("SecondField: got %q, want %q", b.SecondField, "second")
	}
	if wA.Get() != &a {
		t.Error("wA.Get() should return the bound *confA pointer")
	}
	if wA.Viper() == wB.Viper() {
		t.Error("instances should keep separate viper instances")
	}
}

func Test_BindConfig_initErrorReturnedFromExecute(t *testing.T) {
	var testConfig struct{ Field string }

	cmd := &cobra.Command{
		Use:           "test",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          func(*cobra.Command, []string) error { return nil },
	}

	err := BindConfig(&testConfig, cmd, ConfigOpts{
		ConfigName:  "definitely-missing-config",
		ConfigPaths: []string{"."},
		ConfigType:  "json",
		Viper:       globalViper.New(),
	})
	if err != nil {
		t.Fatal(err)
	}

	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected Execute to return an error when the config file is missing, got nil")
	}
}

func Test_BindConfig_chainsPersistentPreRunE(t *testing.T) {
	t.Setenv("CHAIN_FIELD", "from-env")

	var testConfig struct{ ChainField string }
	var sawValue string

	cmd := &cobra.Command{
		Use: "test",
		PersistentPreRunE: func(*cobra.Command, []string) error {
			sawValue = testConfig.ChainField // config must already be initialized here
			return nil
		},
		RunE: func(*cobra.Command, []string) error { return nil },
	}

	if err := BindConfig(&testConfig, cmd, ConfigOpts{Viper: globalViper.New()}); err != nil {
		t.Fatal(err)
	}

	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if sawValue != "from-env" {
		t.Errorf("user PersistentPreRunE saw %q, want %q", sawValue, "from-env")
	}
}

func Test_BindConfig_moreSliceFlagTypes(t *testing.T) {
	var testConfig struct {
		Durations []time.Duration
		Bools     []bool
		Ips       []net.IP
	}

	cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	err := BindConfig(&testConfig, cmd, ConfigOpts{Viper: globalViper.New()})
	if err != nil {
		t.Fatal(err)
	}

	cmd.SetArgs([]string{
		"--durations", "1s,2s",
		"--bools", "true,false",
		"--ips", "10.0.0.1,10.0.0.2",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if len(testConfig.Durations) != 2 || testConfig.Durations[0] != time.Second || testConfig.Durations[1] != 2*time.Second {
		t.Errorf("Durations: got %v, want [1s 2s]", testConfig.Durations)
	}
	if len(testConfig.Bools) != 2 || testConfig.Bools[0] != true || testConfig.Bools[1] != false {
		t.Errorf("Bools: got %v, want [true false]", testConfig.Bools)
	}
	if len(testConfig.Ips) != 2 || testConfig.Ips[0].String() != "10.0.0.1" || testConfig.Ips[1].String() != "10.0.0.2" {
		t.Errorf("Ips: got %v, want [10.0.0.1 10.0.0.2]", testConfig.Ips)
	}
}

func Test_BindConfig_allFlagTypesRegistered(t *testing.T) {
	var testConfig struct {
		S    string
		B    bool
		I    int
		I8   int8
		I16  int16
		I32  int32
		I64  int64
		U    uint
		U8   uint8
		U16  uint16
		U32  uint32
		U64  uint64
		F32  float32
		F64  float64
		D    time.Duration
		T    time.Time
		Ip   net.IP
		Nt   net.IPNet
		Ss   []string
		Is   []int
		I32s []int32
		I64s []int64
		Us   []uint
		Bs   []byte
		F32s []float32
		F64s []float64
		Arr  [2]string
		Ms   map[string]string
		Mi   map[string]int
		Mi32 map[string]int32
		Mi64 map[string]int64
	}

	cmd := &cobra.Command{}
	if err := BindConfig(&testConfig, cmd, ConfigOpts{Viper: globalViper.New()}); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{
		"s", "b", "i", "i-8", "i-16", "i-32", "i-64",
		"u", "u-8", "u-16", "u-32", "u-64", "f-32", "f-64",
		"d", "t", "ip", "nt", "ss", "is", "i-32-s", "i-64-s",
		"us", "bs", "f-32-s", "f-64-s", "arr", "ms", "mi", "mi-32", "mi-64",
	} {
		if cmd.PersistentFlags().Lookup(name) == nil {
			t.Errorf("flag %q not registered", name)
		}
	}
}

func TestStructDefaultsSurviveFlagBinding(t *testing.T) {
	type conf struct {
		Name  string
		Count int
	}
	c := conf{Name: "struct-default", Count: 7}

	cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	if err := BindConfig(&c, cmd, ConfigOpts{Viper: globalViper.New()}); err != nil {
		t.Fatal(err)
	}
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if c.Name != "struct-default" {
		t.Errorf("Name: got %q, want %q (flag zero value must not clobber struct default)", c.Name, "struct-default")
	}
	if c.Count != 7 {
		t.Errorf("Count: got %d, want 7", c.Count)
	}
}

func TestStructDefaultsOverriddenBySources(t *testing.T) {
	t.Setenv("COUNT", "11")

	type conf struct {
		Name  string
		Count int
	}
	c := conf{Name: "struct-default", Count: 7}

	cmd := &cobra.Command{Run: func(*cobra.Command, []string) {}}
	if err := BindConfig(&c, cmd, ConfigOpts{Viper: globalViper.New()}); err != nil {
		t.Fatal(err)
	}
	cmd.SetArgs([]string{"--name", "from-flag"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if c.Name != "from-flag" {
		t.Errorf("Name: got %q, want %q (flag must override struct default)", c.Name, "from-flag")
	}
	if c.Count != 11 {
		t.Errorf("Count: got %d, want 11 (env must override struct default)", c.Count)
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
