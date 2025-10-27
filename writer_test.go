package stdcli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ddollar/errors"
)

func TestWriterWrite(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		color   bool
		want    string
	}{
		{
			name:  "plain text",
			input: "hello world",
			color: false,
			want:  "hello world",
		},
		{
			name:  "text with h1 tag without color",
			input: "<h1>Title</h1>",
			color: false,
			want:  "Title",
		},
		{
			name:  "text with value tag without color",
			input: "<value>test</value>",
			color: false,
			want:  "test",
		},
		{
			name:  "multiple tags without color",
			input: "<h1>Title</h1> <value>value</value>",
			color: false,
			want:  "Title value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &Writer{
				Stdout: buf,
				Stderr: buf,
				Color:  tt.color,
				Tags:   DefaultWriter.Tags,
			}

			n, err := w.Write([]byte(tt.input))
			if err != nil {
				t.Errorf("Write() error = %v", err)
				return
			}

			if n == 0 {
				t.Errorf("Write() n = 0, want > 0")
			}

			got := stripColor(buf.String())
			if got != tt.want {
				t.Errorf("Write() output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriterWritef(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []any
		color  bool
		want   string
	}{
		{
			name:   "simple format",
			format: "hello %s",
			args:   []any{"world"},
			color:  false,
			want:   "hello world",
		},
		{
			name:   "format with tag",
			format: "<h1>%s</h1>",
			args:   []any{"Title"},
			color:  false,
			want:   "Title",
		},
		{
			name:   "format with multiple args",
			format: "<h1>%s</h1> <value>%d</value>",
			args:   []any{"Count", 42},
			color:  false,
			want:   "Count 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &Writer{
				Stdout: buf,
				Stderr: buf,
				Color:  tt.color,
				Tags:   DefaultWriter.Tags,
			}

			n, err := w.Writef(tt.format, tt.args...)
			if err != nil {
				t.Errorf("Writef() error = %v", err)
				return
			}

			if n == 0 {
				t.Errorf("Writef() n = 0, want > 0")
			}

			got := stripColor(buf.String())
			if got != tt.want {
				t.Errorf("Writef() output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriterError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantMsg   string
	}{
		{
			name:    "simple error",
			err:     errors.Errorf("test error"),
			wantMsg: "ERROR: test error",
		},
		{
			name:    "error with formatting",
			err:     errors.Errorf("error: %s", "failed"),
			wantMsg: "ERROR: error: failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &Writer{
				Stdout: buf,
				Stderr: buf,
				Color:  false,
				Tags:   DefaultWriter.Tags,
			}

			err := w.Error(tt.err)
			if err == nil {
				t.Errorf("Error() returned nil, want error")
				return
			}

			got := stripColor(buf.String())
			if !strings.Contains(got, tt.wantMsg) {
				t.Errorf("Error() output = %q, want it to contain %q", got, tt.wantMsg)
			}
		})
	}
}

func TestWriterErrorf(t *testing.T) {
	buf := &bytes.Buffer{}
	w := &Writer{
		Stdout: buf,
		Stderr: buf,
		Color:  false,
		Tags:   DefaultWriter.Tags,
	}

	err := w.Errorf("test %s", "error")
	if err == nil {
		t.Errorf("Errorf() returned nil, want error")
		return
	}

	got := stripColor(buf.String())
	if !strings.Contains(got, "ERROR: test error") {
		t.Errorf("Errorf() output = %q, want it to contain 'ERROR: test error'", got)
	}
}

func TestWriterSprintf(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []any
		want   string
	}{
		{
			name:   "simple format",
			format: "hello %s",
			args:   []any{"world"},
			want:   "hello world",
		},
		{
			name:   "format with tags",
			format: "<h1>%s</h1> <value>%d</value>",
			args:   []any{"Title", 42},
			want:   "Title 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				Color: false,
				Tags:  DefaultWriter.Tags,
			}

			got := stripColor(w.Sprintf(tt.format, tt.args...))
			if got != tt.want {
				t.Errorf("Sprintf() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no color codes",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "with color codes",
			input: "\033[38;5;244mhello\033[0m world",
			want:  "hello world",
		},
		{
			name:  "multiple color codes",
			input: "\033[38;5;244mhello\033[0m \033[38;5;251mworld\033[0m",
			want:  "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripColor(tt.input)
			if got != tt.want {
				t.Errorf("stripColor() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripTag(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name:  "no tags",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "h1 tag",
			input: "<h1>hello</h1>",
			want:  "hello",
		},
		{
			name:  "value tag",
			input: "<value>world</value>",
			want:  "world",
		},
		{
			name:  "nested content",
			input: "<h1>test value</h1>",
			want:  "test value",
		},
		{
			name:  "non-string input",
			input: 42,
			want:  "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTag(tt.input)
			if got != tt.want {
				t.Errorf("stripTag() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripTags(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name:  "no tags",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "single tag",
			input: "<h1>hello</h1>",
			want:  "hello",
		},
		{
			name:  "multiple tags",
			input: "<h1>hello</h1> <value>world</value>",
			want:  "hello world",
		},
		{
			name:  "nested tags",
			input: "<h1><value>nested</value></h1>",
			want:  "nested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTags(tt.input)
			if got != tt.want {
				t.Errorf("stripTags() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderColors(t *testing.T) {
	renderer := RenderColors(244)
	input := "<h1>test</h1>"
	output := renderer(input)

	// Should contain color codes
	if !strings.Contains(output, "\033[") {
		t.Errorf("RenderColors() output should contain color codes")
	}

	// Should contain reset code
	if !strings.Contains(output, "\033[0m") {
		t.Errorf("RenderColors() output should contain reset code")
	}

	// Should contain the text
	if !strings.Contains(output, "test") {
		t.Errorf("RenderColors() output should contain 'test'")
	}
}

func TestRenderUnderline(t *testing.T) {
	renderer := RenderUnderline()
	input := "<u>test</u>"
	output := renderer(input)

	// Should contain underline codes
	if !strings.Contains(output, "\033[4m") {
		t.Errorf("RenderUnderline() output should contain underline start code")
	}

	if !strings.Contains(output, "\033[24m") {
		t.Errorf("RenderUnderline() output should contain underline end code")
	}

	// Should contain the text
	if !strings.Contains(output, "test") {
		t.Errorf("RenderUnderline() output should contain 'test'")
	}
}
