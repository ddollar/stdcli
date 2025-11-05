package stdcli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestContextArg(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		index int
		want  string
	}{
		{
			name:  "first arg",
			args:  []string{"arg1", "arg2", "arg3"},
			index: 0,
			want:  "arg1",
		},
		{
			name:  "second arg",
			args:  []string{"arg1", "arg2", "arg3"},
			index: 1,
			want:  "arg2",
		},
		{
			name:  "last arg",
			args:  []string{"arg1", "arg2", "arg3"},
			index: 2,
			want:  "arg3",
		},
		{
			name:  "out of bounds returns empty",
			args:  []string{"arg1", "arg2"},
			index: 5,
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &defaultContext{
				Context: context.Background(),
				args:    tt.args,
			}

			got := ctx.Arg(tt.index)
			if got != tt.want {
				t.Errorf("Arg(%d) = %q, want %q", tt.index, got, tt.want)
			}
		})
	}
}

func TestContextArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "multiple args",
			args: []string{"arg1", "arg2", "arg3"},
		},
		{
			name: "single arg",
			args: []string{"arg1"},
		},
		{
			name: "no args",
			args: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &defaultContext{
				Context: context.Background(),
				args:    tt.args,
			}

			got := ctx.Args()

			if len(got) != len(tt.args) {
				t.Fatalf("Args() length = %d, want %d", len(got), len(tt.args))
			}

			for i, arg := range tt.args {
				if got[i] != arg {
					t.Errorf("Args()[%d] = %q, want %q", i, got[i], arg)
				}
			}
		})
	}
}

func TestContextFlags(t *testing.T) {
	verboseFlag := BoolFlag("verbose", "v", "verbose")
	verboseFlag.Value = true

	nameFlag := StringFlag("name", "n", "name")
	nameFlag.Value = "test"

	flags := Flags{&verboseFlag, &nameFlag}

	ctx := &defaultContext{
		Context: context.Background(),
		flags:   flags,
	}

	gotFlags := ctx.Flags()

	if len(gotFlags) != 2 {
		t.Fatalf("Flags() length = %d, want 2", len(gotFlags))
	}

	if !gotFlags.Bool("verbose") {
		t.Errorf("verbose flag should be true")
	}

	if gotFlags.String("name") != "test" {
		t.Errorf("name flag = %q, want %q", gotFlags.String("name"), "test")
	}
}

func TestContextWrite(t *testing.T) {
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

	data := []byte("hello world")
	n, err := ctx.Write(data)

	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	if n == 0 {
		t.Errorf("Write() n = 0, want > 0")
	}

	if buf.String() != string(data) {
		t.Errorf("Write() output = %q, want %q", buf.String(), data)
	}
}

func TestContextWritef(t *testing.T) {
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

	ctx.Writef("hello %s", "world")

	output := buf.String()
	if output != "hello world" {
		t.Errorf("Writef() output = %q, want %q", output, "hello world")
	}
}

func TestContextVersion(t *testing.T) {
	ctx := &defaultContext{
		Context: context.Background(),
		engine: &Engine{
			Version: "1.2.3",
		},
	}

	got := ctx.Version()
	if got != "1.2.3" {
		t.Errorf("Version() = %q, want %q", got, "1.2.3")
	}
}

func TestContextInfo(t *testing.T) {
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
	if info == nil {
		t.Fatal("Info() returned nil")
	}

	// Test that we can use the InfoWriter
	info.Add("Key", "Value")
	err := info.Print()
	if err != nil {
		t.Errorf("InfoWriter.Print() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "KEY") || !strings.Contains(output, "Value") {
		t.Errorf("InfoWriter output should contain KEY and Value, got: %q", output)
	}
}

func TestContextTable(t *testing.T) {
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

	table := ctx.Table("ID", "Name")
	if table == nil {
		t.Fatal("Table() returned nil")
	}

	// Test that we can use the TableWriter
	table.Append(1, "test")
	err := table.Print()
	if err != nil {
		t.Errorf("TableWriter.Print() error = %v", err)
	}

	output := stripTags(buf.String())
	if !strings.Contains(output, "ID") || !strings.Contains(output, "Name") || !strings.Contains(output, "test") {
		t.Errorf("TableWriter output should contain ID, Name, and test, got: %q", output)
	}
}

func TestContextColumns(t *testing.T) {
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

	columns := ctx.Columns()
	if columns == nil {
		t.Fatal("Columns() returned nil")
	}

	// Test that we can use the ColumnWriter
	columns.Append("col1", "col2", "col3")
	err := columns.Print()
	if err != nil {
		t.Errorf("ColumnWriter.Print() error = %v", err)
	}

	output := stripTags(buf.String())
	if !strings.Contains(output, "col1") || !strings.Contains(output, "col2") || !strings.Contains(output, "col3") {
		t.Errorf("ColumnWriter output should contain col1, col2, and col3, got: %q", output)
	}
}

func TestContextCleanup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Use a channel for proper synchronization
	cleanupDone := make(chan bool, 1)

	dc := &defaultContext{
		Context: ctx,
		engine:  &Engine{},
	}

	dc.Cleanup(func() {
		cleanupDone <- true
	})

	// Cancel the context to trigger cleanup
	cancel()

	// Wait for cleanup to complete with timeout
	select {
	case <-cleanupDone:
		// Success
	case <-context.Background().Done():
		t.Errorf("cleanup function was not called after context cancellation")
	}
}

func TestContextIsTerminal(t *testing.T) {
	// Test with non-terminal reader/writer
	buf := &bytes.Buffer{}

	ctx := &defaultContext{
		Context: context.Background(),
		engine: &Engine{
			Reader: &Reader{Reader: buf},
			Writer: &Writer{Stdout: buf, Stderr: buf},
		},
	}

	// Buffer is not a terminal
	if ctx.IsTerminal() {
		t.Errorf("IsTerminal() = true for buffer, want false")
	}

	if ctx.IsTerminalReader() {
		t.Errorf("IsTerminalReader() = true for buffer, want false")
	}

	if ctx.IsTerminalWriter() {
		t.Errorf("IsTerminalWriter() = true for buffer, want false")
	}
}

func TestContextRead(t *testing.T) {
	input := bytes.NewBufferString("hello world")

	ctx := &defaultContext{
		Context: context.Background(),
		engine: &Engine{
			Reader: &Reader{Reader: input},
		},
	}

	data := make([]byte, 5)
	n, err := ctx.Read(data)

	if err != nil {
		t.Errorf("Read() error = %v", err)
	}

	if n != 5 {
		t.Errorf("Read() n = %d, want 5", n)
	}

	if string(data) != "hello" {
		t.Errorf("Read() data = %q, want %q", data, "hello")
	}
}
