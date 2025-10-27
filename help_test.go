package stdcli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestCommandHelp(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedInHelp []string
	}{
		{
			name: "help with --help flag",
			args: []string{"--help"},
			expectedInHelp: []string{
				"USAGE",
				"DESCRIPTION",
				"test command description",
				"OPTIONS",
			},
		},
		{
			name: "help with -h flag",
			args: []string{"-h"},
			expectedInHelp: []string{
				"USAGE",
				"DESCRIPTION",
				"test command description",
				"OPTIONS",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			buf := &bytes.Buffer{}
			writer := &Writer{
				Stdout: buf,
				Stderr: buf,
				Color:  false,
				Tags:   map[string]Renderer{},
			}

			e := &Engine{
				Name:    "testapp",
				Version: "1.0.0",
				Writer:  writer,
			}

			e.Command("test", "test command description", func(ctx Context) error {
				return nil
			}, CommandOptions{
				Usage: "<arg>",
				Flags: []Flag{
					StringFlag("option", "o", "test option"),
					BoolFlag("verbose", "v", "verbose output"),
				},
			})

			// Find the test command
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

			err := cmd.ExecuteContext(context.Background(), tt.args)

			// Help should return nil (not an error)
			if err != nil {
				t.Errorf("help should not return error, got: %v", err)
			}

			output := buf.String()

			// Check for expected content in help
			for _, expected := range tt.expectedInHelp {
				if !strings.Contains(output, expected) {
					t.Errorf("expected help to contain %q, but it didn't. Output:\n%s", expected, output)
				}
			}

			// Should NOT contain pflag's default format
			if strings.Contains(output, "Usage of :") {
				t.Errorf("help should not contain pflag default format. Output:\n%s", output)
			}
		})
	}
}

func TestCommandHelpWithGlobalFlags(t *testing.T) {
	// Capture output
	buf := &bytes.Buffer{}
	writer := &Writer{
		Stdout: buf,
		Stderr: buf,
		Color:  false,
		Tags:   map[string]Renderer{},
	}

	e := &Engine{
		Name:    "testapp",
		Version: "1.0.0",
		Writer:  writer,
		Flags: []Flag{
			BoolFlag("debug", "d", "enable debug mode"),
		},
	}

	e.Command("test", "test command description", func(ctx Context) error {
		return nil
	}, CommandOptions{
		Usage: "<arg>",
		Flags: []Flag{
			StringFlag("option", "o", "test option"),
		},
	})

	// Find the test command
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

	err := cmd.ExecuteContext(context.Background(), []string{"--help"})

	if err != nil {
		t.Errorf("help should not return error, got: %v", err)
	}

	output := buf.String()

	// Should have both OPTIONS and GLOBAL OPTIONS sections
	expectedSections := []string{
		"USAGE",
		"DESCRIPTION",
		"OPTIONS",
		"GLOBAL OPTIONS",
		"--option",
		"--debug",
	}

	for _, expected := range expectedSections {
		if !strings.Contains(output, expected) {
			t.Errorf("expected help to contain %q. Output:\n%s", expected, output)
		}
	}
}
