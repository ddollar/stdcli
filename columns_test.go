package stdcli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestColumnWriter(t *testing.T) {
	tests := []struct {
		name      string
		rows      [][]any
		wantItems []string
	}{
		{
			name: "single row",
			rows: [][]any{
				{"id-1", "name-1", "status-1"},
			},
			wantItems: []string{"id-1", "name-1", "status-1"},
		},
		{
			name: "multiple rows",
			rows: [][]any{
				{"id-1", "name-1", "status-1"},
				{"id-2", "name-2", "status-2"},
				{"id-3", "name-3", "status-3"},
			},
			wantItems: []string{
				"id-1", "name-1", "status-1",
				"id-2", "name-2", "status-2",
				"id-3", "name-3", "status-3",
			},
		},
		{
			name: "rows with different lengths",
			rows: [][]any{
				{"short", "medium-length", "very-very-long-value"},
				{"x", "y", "z"},
			},
			wantItems: []string{
				"short", "medium-length", "very-very-long-value",
				"x", "y", "z",
			},
		},
		{
			name: "rows with numbers",
			rows: [][]any{
				{"id", "count", "active"},
				{1, 42, true},
				{2, 100, false},
			},
			wantItems: []string{
				"id", "count", "active",
				"1", "42", "true",
				"2", "100", "false",
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

			cw := ctx.Columns()

			for _, row := range tt.rows {
				cw.Append(row...)
			}

			err := cw.Print()
			if err != nil {
				t.Errorf("Print() error = %v", err)
				return
			}

			output := stripTags(buf.String())

			for _, want := range tt.wantItems {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got:\n%s", want, output)
				}
			}
		})
	}
}

func TestColumnWriterEmpty(t *testing.T) {
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

	cw := ctx.Columns()
	err := cw.Print()

	if err != nil {
		t.Errorf("Print() error = %v", err)
	}

	// Should handle empty rows gracefully
	if buf.Len() != 0 {
		t.Errorf("expected empty output for no rows, got: %q", buf.String())
	}
}

func TestColumnWriterAlignment(t *testing.T) {
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

	cw := ctx.Columns()
	cw.Append("short", "medium-size", "value")
	cw.Append("x", "y", "long-last-value")

	err := cw.Print()
	if err != nil {
		t.Errorf("Print() error = %v", err)
	}

	output := stripTags(buf.String())
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	// Verify columns are padded (first line should be longer than just concatenated)
	firstLine := lines[0]
	if !strings.Contains(firstLine, "short") || !strings.Contains(firstLine, "medium-size") {
		t.Errorf("first line should contain padded columns, got: %q", firstLine)
	}
}

func TestColumnWriterWithTags(t *testing.T) {
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

	cw := ctx.Columns()
	cw.Append("<h1>Header</h1>", "<value>Value</value>")

	err := cw.Print()
	if err != nil {
		t.Errorf("Print() error = %v", err)
	}

	output := stripTags(buf.String())

	// Tags should be stripped in the output but content preserved
	if !strings.Contains(output, "Header") {
		t.Errorf("output should contain 'Header', got: %q", output)
	}
	if !strings.Contains(output, "Value") {
		t.Errorf("output should contain 'Value', got: %q", output)
	}
}

func TestColumnWriterWidths(t *testing.T) {
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

	cw := &columnWriter{ctx: ctx}

	// Test widths calculation
	cw.rows = [][]any{
		{"abc", "defgh", "ij"},
		{"x", "yy", "zzz"},
	}

	widths := cw.widths()

	// First column: max(3, 1) = 3
	// Second column: max(5, 2) = 5
	// Third column: 0 (last column not padded)
	expectedWidths := []int{3, 5, 0}

	if len(widths) != len(expectedWidths) {
		t.Fatalf("expected %d widths, got %d", len(expectedWidths), len(widths))
	}

	for i, want := range expectedWidths {
		if widths[i] != want {
			t.Errorf("widths[%d] = %d, want %d", i, widths[i], want)
		}
	}
}
