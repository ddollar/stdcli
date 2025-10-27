package stdcli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestTableWriterText(t *testing.T) {
	tests := []struct {
		name      string
		columns   []any
		rows      [][]any
		wantItems []string
	}{
		{
			name:    "simple table",
			columns: []any{"ID", "Name", "Status"},
			rows: [][]any{
				{"1", "test1", "active"},
				{"2", "test2", "inactive"},
			},
			wantItems: []string{"ID", "Name", "Status", "1", "test1", "active", "2", "test2", "inactive"},
		},
		{
			name:    "table with numbers",
			columns: []any{"ID", "Count", "Active"},
			rows: [][]any{
				{1, 100, true},
				{2, 200, false},
			},
			wantItems: []string{"ID", "Count", "Active", "1", "100", "true", "2", "200", "false"},
		},
		{
			name:      "empty table",
			columns:   []any{"ID", "Name"},
			rows:      [][]any{},
			wantItems: []string{"ID", "Name"},
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

			outputFlag := StringFlag("output", "", "output format")

			ctx := &defaultContext{
				Context: context.Background(),
				flags:   Flags{&outputFlag},
				engine: &Engine{
					Writer: writer,
				},
			}

			table := ctx.Table(tt.columns...)

			for _, row := range tt.rows {
				table.Append(row...)
			}

			err := table.Print()
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

func TestTableWriterJSON(t *testing.T) {
	tests := []struct {
		name    string
		columns []any
		rows    [][]any
		want    []map[string]any
	}{
		{
			name:    "simple table",
			columns: []any{"ID", "Name", "Status"},
			rows: [][]any{
				{"1", "test1", "active"},
				{"2", "test2", "inactive"},
			},
			want: []map[string]any{
				{"id": "1", "name": "test1", "status": "active"},
				{"id": "2", "name": "test2", "status": "inactive"},
			},
		},
		{
			name:    "table with mixed types",
			columns: []any{"ID", "Count", "Active"},
			rows: [][]any{
				{1, 100, true},
				{2, 200, false},
			},
			want: []map[string]any{
				{"id": 1, "count": 100, "active": true},
				{"id": 2, "count": 200, "active": false},
			},
		},
		{
			name:    "empty table",
			columns: []any{"ID", "Name"},
			rows:    [][]any{},
			want:    []map[string]any{},
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

			outputFlag := StringFlag("output", "", "output format")
			outputFlag.Value = "json"

			ctx := &defaultContext{
				Context: context.Background(),
				flags:   Flags{&outputFlag},
				engine: &Engine{
					Writer: writer,
				},
			}

			table := ctx.Table(tt.columns...)

			for _, row := range tt.rows {
				table.Append(row...)
			}

			err := table.Print()
			if err != nil {
				t.Errorf("Print() error = %v", err)
				return
			}

			output := buf.String()

			var got []map[string]any
			if err := json.Unmarshal([]byte(output), &got); err != nil {
				t.Errorf("failed to unmarshal JSON: %v, output: %s", err, output)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("got %d rows, want %d rows", len(got), len(tt.want))
				return
			}

			for i, wantRow := range tt.want {
				gotRow := got[i]
				for key, wantVal := range wantRow {
					gotVal, ok := gotRow[key]
					if !ok {
						t.Errorf("row %d missing key %q", i, key)
						continue
					}

					// Handle number comparison (JSON unmarshals numbers as float64)
					if wantInt, ok := wantVal.(int); ok {
						if gotFloat, ok := gotVal.(float64); ok {
							if int(gotFloat) != wantInt {
								t.Errorf("row %d key %q = %v, want %v", i, key, gotVal, wantVal)
							}
							continue
						}
					}

					if gotVal != wantVal {
						t.Errorf("row %d key %q = %v, want %v", i, key, gotVal, wantVal)
					}
				}
			}
		})
	}
}

func TestTableWriterWidths(t *testing.T) {
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

	tw := &tableWriter{
		ctx:     ctx,
		columns: []any{"Short", "MediumLength", "VeryLongColumnName"},
		rows: [][]any{
			{"a", "bb", "ccc"},
			{"dddd", "e", "f"},
		},
	}

	widths := tw.widths()

	// First column: max(5 (Short), 1 (a), 4 (dddd)) = 5
	// Second column: max(12 (MediumLength), 2 (bb), 1 (e)) = 12
	// Third column: 0 (last column not padded)
	expectedWidths := []int{5, 12, 0}

	if len(widths) != len(expectedWidths) {
		t.Fatalf("expected %d widths, got %d", len(expectedWidths), len(widths))
	}

	for i, want := range expectedWidths {
		if widths[i] != want {
			t.Errorf("widths[%d] = %d, want %d", i, widths[i], want)
		}
	}
}

func TestTableWriterColumnCaseLowering(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := &Writer{
		Stdout: buf,
		Stderr: buf,
		Color:  false,
		Tags:   map[string]Renderer{},
	}

	outputFlag := StringFlag("output", "", "output format")
	outputFlag.Value = "json"

	ctx := &defaultContext{
		Context: context.Background(),
		flags:   Flags{&outputFlag},
		engine: &Engine{
			Writer: writer,
		},
	}

	// Test that column names are lowercased in JSON output
	table := ctx.Table("ID", "UserName", "IsActive")
	table.Append(1, "john", true)

	err := table.Print()
	if err != nil {
		t.Errorf("Print() error = %v", err)
		return
	}

	output := buf.String()

	var got []map[string]any
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Errorf("failed to unmarshal JSON: %v", err)
		return
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 row, got %d", len(got))
	}

	// Check that keys are lowercase
	row := got[0]
	if _, ok := row["id"]; !ok {
		t.Errorf("expected key 'id', got keys: %v", getKeys(row))
	}
	if _, ok := row["username"]; !ok {
		t.Errorf("expected key 'username', got keys: %v", getKeys(row))
	}
	if _, ok := row["isactive"]; !ok {
		t.Errorf("expected key 'isactive', got keys: %v", getKeys(row))
	}
}

func getKeys(m map[string]any) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
