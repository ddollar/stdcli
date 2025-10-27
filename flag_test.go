package stdcli

import (
	"strings"
	"testing"
	"time"
)

func TestBoolFlag(t *testing.T) {
	flag := BoolFlag("verbose", "v", "enable verbose output")

	if flag.Name != "verbose" {
		t.Errorf("Name = %q, want %q", flag.Name, "verbose")
	}
	if flag.Short != "v" {
		t.Errorf("Short = %q, want %q", flag.Short, "v")
	}
	if flag.Description != "enable verbose output" {
		t.Errorf("Description = %q, want %q", flag.Description, "enable verbose output")
	}
	if flag.Kind() != FlagBool {
		t.Errorf("Kind() = %q, want %q", flag.Kind(), FlagBool)
	}
}

func TestStringFlag(t *testing.T) {
	flag := StringFlag("name", "n", "specify name")

	if flag.Name != "name" {
		t.Errorf("Name = %q, want %q", flag.Name, "name")
	}
	if flag.Short != "n" {
		t.Errorf("Short = %q, want %q", flag.Short, "n")
	}
	if flag.Description != "specify name" {
		t.Errorf("Description = %q, want %q", flag.Description, "specify name")
	}
	if flag.Kind() != FlagString {
		t.Errorf("Kind() = %q, want %q", flag.Kind(), FlagString)
	}
}

func TestIntFlag(t *testing.T) {
	flag := IntFlag("count", "c", "specify count")

	if flag.Name != "count" {
		t.Errorf("Name = %q, want %q", flag.Name, "count")
	}
	if flag.Short != "c" {
		t.Errorf("Short = %q, want %q", flag.Short, "c")
	}
	if flag.Description != "specify count" {
		t.Errorf("Description = %q, want %q", flag.Description, "specify count")
	}
	if flag.Kind() != FlagInt {
		t.Errorf("Kind() = %q, want %q", flag.Kind(), FlagInt)
	}
}

func TestDurationFlag(t *testing.T) {
	flag := DurationFlag("timeout", "t", "specify timeout")

	if flag.Name != "timeout" {
		t.Errorf("Name = %q, want %q", flag.Name, "timeout")
	}
	if flag.Short != "t" {
		t.Errorf("Short = %q, want %q", flag.Short, "t")
	}
	if flag.Description != "specify timeout" {
		t.Errorf("Description = %q, want %q", flag.Description, "specify timeout")
	}
	if flag.Kind() != FlagDuration {
		t.Errorf("Kind() = %q, want %q", flag.Kind(), FlagDuration)
	}
}

func TestFlagSet(t *testing.T) {
	tests := []struct {
		name      string
		flag      Flag
		setValue  string
		wantValue any
		wantError bool
	}{
		{
			name:      "set bool flag to true",
			flag:      BoolFlag("test", "", "test flag"),
			setValue:  "true",
			wantValue: true,
			wantError: false,
		},
		{
			name:      "set bool flag to false",
			flag:      BoolFlag("test", "", "test flag"),
			setValue:  "false",
			wantValue: false,
			wantError: false,
		},
		{
			name:      "set string flag",
			flag:      StringFlag("test", "", "test flag"),
			setValue:  "hello",
			wantValue: "hello",
			wantError: false,
		},
		{
			name:      "set int flag",
			flag:      IntFlag("test", "", "test flag"),
			setValue:  "42",
			wantValue: 42,
			wantError: false,
		},
		{
			name:      "set int flag with invalid value",
			flag:      IntFlag("test", "", "test flag"),
			setValue:  "notanumber",
			wantError: true,
		},
		{
			name:      "set duration flag",
			flag:      DurationFlag("test", "", "test flag"),
			setValue:  "5s",
			wantValue: 5 * time.Second,
			wantError: false,
		},
		{
			name:      "set duration flag with invalid value",
			flag:      DurationFlag("test", "", "test flag"),
			setValue:  "invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.flag.Set(tt.setValue)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.flag.Value != tt.wantValue {
				t.Errorf("Value = %v, want %v", tt.flag.Value, tt.wantValue)
			}
		})
	}
}

func TestFlagUsage(t *testing.T) {
	tests := []struct {
		name  string
		flag  Flag
		want  string
	}{
		{
			name: "bool flag with short",
			flag: BoolFlag("verbose", "v", "enable verbose"),
			want: "-v --verbose",
		},
		{
			name: "bool flag without short",
			flag: BoolFlag("verbose", "", "enable verbose"),
			want: "   --verbose",
		},
		{
			name: "string flag with short",
			flag: StringFlag("name", "n", "specify name"),
			want: "-n --name <name>",
		},
		{
			name: "string flag without short",
			flag: StringFlag("name", "", "specify name"),
			want: "   --name <name>",
		},
		{
			name: "int flag with short",
			flag: IntFlag("count", "c", "specify count"),
			want: "-c --count <count>",
		},
		{
			name: "duration flag with short",
			flag: DurationFlag("timeout", "t", "specify timeout"),
			want: "-t --timeout <timeout>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTags(tt.flag.Usage())
			want := tt.want

			if got != want {
				t.Errorf("Usage() = %q, want %q", got, want)
			}
		})
	}
}

func TestFlagsAccessors(t *testing.T) {
	boolFlag := BoolFlag("bool-flag", "", "bool flag")
	boolFlag.Default = true
	boolFlag.Value = false

	stringFlag := StringFlag("string-flag", "", "string flag")
	stringFlag.Default = "default"
	stringFlag.Value = "value"

	intFlag := IntFlag("int-flag", "", "int flag")
	intFlag.Default = 10
	intFlag.Value = 42

	flags := Flags{&boolFlag, &stringFlag, &intFlag}

	// Test Bool accessor
	if got := flags.Bool("bool-flag"); got != false {
		t.Errorf("Bool() = %v, want false", got)
	}

	// Test Bool accessor with default
	if got := flags.Bool("nonexistent"); got != false {
		t.Errorf("Bool() for nonexistent = %v, want false", got)
	}

	// Test String accessor
	if got := flags.String("string-flag"); got != "value" {
		t.Errorf("String() = %q, want %q", got, "value")
	}

	// Test String accessor with default
	if got := flags.String("nonexistent"); got != "" {
		t.Errorf("String() for nonexistent = %q, want empty string", got)
	}

	// Test Int accessor
	if got := flags.Int("int-flag"); got != 42 {
		t.Errorf("Int() = %d, want 42", got)
	}

	// Test Int accessor with default
	if got := flags.Int("nonexistent"); got != 0 {
		t.Errorf("Int() for nonexistent = %d, want 0", got)
	}

	// Test Value accessor
	if got := flags.Value("string-flag"); got != "value" {
		t.Errorf("Value() = %v, want %q", got, "value")
	}

	// Test Value accessor for nonexistent
	if got := flags.Value("nonexistent"); got != nil {
		t.Errorf("Value() for nonexistent = %v, want nil", got)
	}
}

func TestFlagsWithNilValues(t *testing.T) {
	boolFlag := BoolFlag("bool-flag", "", "bool flag")
	boolFlag.Default = true

	stringFlag := StringFlag("string-flag", "", "string flag")
	stringFlag.Default = "default"

	intFlag := IntFlag("int-flag", "", "int flag")
	intFlag.Default = 10

	flags := Flags{&boolFlag, &stringFlag, &intFlag}

	// When Value is nil, should return default
	if got := flags.Bool("bool-flag"); got != true {
		t.Errorf("Bool() with nil value = %v, want true (default)", got)
	}

	if got := flags.String("string-flag"); got != "default" {
		t.Errorf("String() with nil value = %q, want %q (default)", got, "default")
	}

	if got := flags.Int("int-flag"); got != 10 {
		t.Errorf("Int() with nil value = %d, want 10 (default)", got)
	}
}

func TestFlagString(t *testing.T) {
	flag := StringFlag("test", "", "test flag")
	flag.Value = "hello"

	got := flag.String()
	if !strings.Contains(got, "test") || !strings.Contains(got, "hello") {
		t.Errorf("String() = %q, want it to contain both 'test' and 'hello'", got)
	}
}

func TestFlagType(t *testing.T) {
	tests := []struct {
		name     string
		flag     Flag
		wantType string
	}{
		{
			name:     "bool flag",
			flag:     BoolFlag("test", "", "test"),
			wantType: "bool",
		},
		{
			name:     "string flag",
			flag:     StringFlag("test", "", "test"),
			wantType: "string",
		},
		{
			name:     "int flag",
			flag:     IntFlag("test", "", "test"),
			wantType: "int",
		},
		{
			name:     "duration flag",
			flag:     DurationFlag("test", "", "test"),
			wantType: "duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flag.Type(); got != tt.wantType {
				t.Errorf("Type() = %q, want %q", got, tt.wantType)
			}
		})
	}
}
