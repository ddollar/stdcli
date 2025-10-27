package stdcli

import (
	"bytes"
	"context"
	"testing"
)

func TestEngineNew(t *testing.T) {
	e := New("testapp", "1.0.0")

	if e.Name != "testapp" {
		t.Errorf("Name = %q, want %q", e.Name, "testapp")
	}

	if e.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", e.Version, "1.0.0")
	}

	// Should have default help command
	if len(e.Commands) == 0 {
		t.Errorf("expected at least one command (help)")
	}

	helpFound := false
	for _, cmd := range e.Commands {
		if len(cmd.Command) > 0 && cmd.Command[0] == "help" {
			helpFound = true
			break
		}
	}

	if !helpFound {
		t.Errorf("expected default 'help' command")
	}
}

func TestEngineCommand(t *testing.T) {
	e := New("testapp", "1.0.0")

	handlerCalled := false
	e.Command("test", "test command", func(ctx Context) error {
		handlerCalled = true
		return nil
	}, CommandOptions{})

	if len(e.Commands) < 2 {
		t.Fatalf("expected at least 2 commands, got %d", len(e.Commands))
	}

	// Find the test command
	var testCmd *Command
	for i, cmd := range e.Commands {
		if len(cmd.Command) > 0 && cmd.Command[0] == "test" {
			testCmd = &e.Commands[i]
			break
		}
	}

	if testCmd == nil {
		t.Fatal("test command not found")
	}

	if testCmd.Description != "test command" {
		t.Errorf("Description = %q, want %q", testCmd.Description, "test command")
	}

	// Execute the command
	err := testCmd.ExecuteContext(context.Background(), []string{})
	if err != nil {
		t.Errorf("ExecuteContext() error = %v", err)
	}

	if !handlerCalled {
		t.Errorf("handler was not called")
	}
}

func TestEngineMultiWordCommand(t *testing.T) {
	e := New("testapp", "1.0.0")

	handlerCalled := false
	e.Command("resource list", "list resources", func(ctx Context) error {
		handlerCalled = true
		return nil
	}, CommandOptions{})

	// Find the command
	var cmd *Command
	for i, c := range e.Commands {
		if len(c.Command) == 2 && c.Command[0] == "resource" && c.Command[1] == "list" {
			cmd = &e.Commands[i]
			break
		}
	}

	if cmd == nil {
		t.Fatal("multi-word command not found")
	}

	if len(cmd.Command) != 2 {
		t.Errorf("expected 2-part command, got %d parts", len(cmd.Command))
	}

	// Execute with matching args
	err := cmd.ExecuteContext(context.Background(), []string{})
	if err != nil {
		t.Errorf("ExecuteContext() error = %v", err)
	}

	if !handlerCalled {
		t.Errorf("handler was not called")
	}
}

func TestEngineExecuteContext(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantExitCode   int
		setupCommands  func(*Engine)
	}{
		{
			name:         "successful command",
			args:         []string{"test"},
			wantExitCode: 0,
			setupCommands: func(e *Engine) {
				e.Command("test", "test command", func(ctx Context) error {
					return nil
				}, CommandOptions{})
			},
		},
		{
			name:         "command with error",
			args:         []string{"test"},
			wantExitCode: 1,
			setupCommands: func(e *Engine) {
				e.Command("test", "test command", func(ctx Context) error {
					return Exit(1).(error)
				}, CommandOptions{})
			},
		},
		{
			name:         "command with custom exit code",
			args:         []string{"test"},
			wantExitCode: 42,
			setupCommands: func(e *Engine) {
				e.Command("test", "test command", func(ctx Context) error {
					return Exit(42).(error)
				}, CommandOptions{})
			},
		},
		{
			name:         "version flag",
			args:         []string{"--version"},
			wantExitCode: 0,
			setupCommands: func(e *Engine) {
				// No additional commands needed
			},
		},
		{
			name:         "version short flag",
			args:         []string{"-v"},
			wantExitCode: 0,
			setupCommands: func(e *Engine) {
				// No additional commands needed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			e := &Engine{
				Name:    "testapp",
				Version: "1.0.0",
				Writer: &Writer{
					Stdout: buf,
					Stderr: buf,
					Color:  false,
					Tags:   DefaultWriter.Tags,
				},
			}

			// Add default help command
			e.Command("help", "show help", func(ctx Context) error {
				return nil
			}, CommandOptions{})

			tt.setupCommands(e)

			exitCode := e.ExecuteContext(context.Background(), tt.args)

			if exitCode != tt.wantExitCode {
				t.Errorf("ExecuteContext() exitCode = %d, want %d, output: %s", exitCode, tt.wantExitCode, buf.String())
			}
		})
	}
}

func TestEngineCommandMatching(t *testing.T) {
	e := New("testapp", "1.0.0")

	var calledCommand string

	e.Command("short", "short command", func(ctx Context) error {
		calledCommand = "short"
		return nil
	}, CommandOptions{})

	e.Command("short long", "long command", func(ctx Context) error {
		calledCommand = "short long"
		return nil
	}, CommandOptions{})

	tests := []struct {
		name        string
		args        []string
		wantCommand string
	}{
		{
			name:        "matches short command",
			args:        []string{"short"},
			wantCommand: "short",
		},
		{
			name:        "matches longer command",
			args:        []string{"short", "long"},
			wantCommand: "short long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calledCommand = ""

			// Find the matching command
			var matchedCmd *Command
			var matchedArgs []string

			for _, c := range e.Commands {
				if args, ok := c.Match(tt.args); ok {
					if matchedCmd == nil || len(matchedCmd.Command) < len(c.Command) {
						cmd := c
						matchedCmd = &cmd
						matchedArgs = args
					}
				}
			}

			if matchedCmd == nil {
				t.Fatal("no command matched")
			}

			err := matchedCmd.ExecuteContext(context.Background(), matchedArgs)
			if err != nil {
				t.Errorf("ExecuteContext() error = %v", err)
			}

			if calledCommand != tt.wantCommand {
				t.Errorf("called command = %q, want %q", calledCommand, tt.wantCommand)
			}
		})
	}
}

func TestEngineGlobalFlags(t *testing.T) {
	buf := &bytes.Buffer{}
	e := &Engine{
		Name:    "testapp",
		Version: "1.0.0",
		Writer: &Writer{
			Stdout: buf,
			Stderr: buf,
			Color:  false,
			Tags:   map[string]Renderer{},
		},
		Flags: []Flag{
			BoolFlag("debug", "d", "enable debug"),
		},
	}

	debugFlagValue := false

	e.Command("test", "test command", func(ctx Context) error {
		debugFlagValue = ctx.Flags().Bool("debug")
		return nil
	}, CommandOptions{})

	// Find test command
	var cmd *Command
	for i, c := range e.Commands {
		if len(c.Command) > 0 && c.Command[0] == "test" {
			cmd = &e.Commands[i]
			break
		}
	}

	if cmd == nil {
		t.Fatal("test command not found")
	}

	// Execute with global flag
	err := cmd.ExecuteContext(context.Background(), []string{"--debug"})
	if err != nil {
		t.Errorf("ExecuteContext() error = %v", err)
	}

	if !debugFlagValue {
		t.Errorf("debug flag should be true")
	}
}

func TestEngineInvisibleCommand(t *testing.T) {
	e := New("testapp", "1.0.0")

	e.Command("visible", "visible command", func(ctx Context) error {
		return nil
	}, CommandOptions{})

	e.Command("invisible", "invisible command", func(ctx Context) error {
		return nil
	}, CommandOptions{
		Invisible: true,
	})

	visibleCount := 0
	for _, cmd := range e.Commands {
		if !cmd.Invisible {
			visibleCount++
		}
	}

	// Should have help and visible commands, but not invisible
	if visibleCount != 2 {
		t.Errorf("expected 2 visible commands (help + visible), got %d", visibleCount)
	}
}

func TestEngineVersionExitCode(t *testing.T) {
	e := &Engine{
		Name:    "testapp",
		Version: "2.5.3",
		Writer: &Writer{
			Stdout: &bytes.Buffer{},
			Stderr: &bytes.Buffer{},
			Color:  false,
			Tags:   map[string]Renderer{},
		},
	}

	// Test that --version returns exit code 0
	exitCode := e.ExecuteContext(context.Background(), []string{"--version"})
	if exitCode != 0 {
		t.Errorf("version command should exit with 0, got %d", exitCode)
	}

	// Test that -v returns exit code 0
	exitCode = e.ExecuteContext(context.Background(), []string{"-v"})
	if exitCode != 0 {
		t.Errorf("version short flag should exit with 0, got %d", exitCode)
	}

	// Note: Version output uses fmt.Println which writes directly to os.Stdout,
	// not to the engine's Writer, so we can't capture it in tests
}
