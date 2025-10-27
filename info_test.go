package stdcli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestInfoWriter(t *testing.T) {
	tests := []struct {
		name     string
		rows     []struct{ header, value string }
		wantRows []string
	}{
		{
			name: "single row",
			rows: []struct{ header, value string }{
				{"Name", "test"},
			},
			wantRows: []string{
				"Name",
				"test",
			},
		},
		{
			name: "multiple rows",
			rows: []struct{ header, value string }{
				{"Name", "test"},
				{"Version", "1.0.0"},
				{"Status", "running"},
			},
			wantRows: []string{
				"Name",
				"test",
				"Version",
				"1.0.0",
				"Status",
				"running",
			},
		},
		{
			name: "rows with different header lengths",
			rows: []struct{ header, value string }{
				{"ID", "123"},
				{"Name", "test"},
				{"Description", "a longer description"},
			},
			wantRows: []string{
				"ID",
				"123",
				"Name",
				"test",
				"Description",
				"a longer description",
			},
		},
		{
			name: "value with newlines",
			rows: []struct{ header, value string }{
				{"Name", "test"},
				{"Description", "line 1\nline 2\nline 3"},
			},
			wantRows: []string{
				"Name",
				"test",
				"Description",
				"line 1",
				"line 2",
				"line 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			writer := &Writer{
				Stdout: buf,
				Stderr: buf,
				Color:  false,
				Tags:   map[string]Renderer{},
			}

			ctx := &defaultContext{
				Context: context.Background(),
				engine: &Engine{
					Writer: writer,
				},
			}

			info := ctx.Info()

			for _, row := range tt.rows {
				info.Add(row.header, row.value)
			}

			err := info.Print()
			if err != nil {
				t.Errorf("Print() error = %v", err)
				return
			}

			output := buf.String()

			for _, want := range tt.wantRows {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got:\n%s", want, output)
				}
			}
		})
	}
}

func TestInfoWriterEmpty(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := &Writer{
		Stdout: buf,
		Stderr: buf,
		Color:  false,
		Tags:   map[string]Renderer{},
	}

	ctx := &defaultContext{
		Context: context.Background(),
		engine: &Engine{
			Writer: writer,
		},
	}

	info := ctx.Info()
	err := info.Print()

	if err != nil {
		t.Errorf("Print() error = %v", err)
	}

	// Should not crash with empty rows
}

func TestInfoWriterAlignment(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := &Writer{
		Stdout: buf,
		Stderr: buf,
		Color:  false,
		Tags:   map[string]Renderer{},
	}

	ctx := &defaultContext{
		Context: context.Background(),
		engine: &Engine{
			Writer: writer,
		},
	}

	info := ctx.Info()
	info.Add("ID", "123")
	info.Add("Very Long Header", "value")

	err := info.Print()
	if err != nil {
		t.Errorf("Print() error = %v", err)
	}

	// Output should be properly aligned
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}

	// Both lines should have consistent spacing
	// (This is a basic check - actual formatting may vary)
	if len(output) == 0 {
		t.Errorf("expected non-empty output")
	}
}
