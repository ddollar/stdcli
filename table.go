package stdcli

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TableWriter interface {
	Append(row ...any)
	Print() error
}

type tableWriter struct {
	ctx     Context
	columns []any
	rows    [][]any
}

var _ TableWriter = &tableWriter{}

func (t *tableWriter) Append(row ...any) {
	t.rows = append(t.rows, row)
}

func (t *tableWriter) Print() error {
	switch t.ctx.Flags().String("output") {
	case "json":
		return t.printJSON()
	default:
		return t.printText()
	}
}

func (t *tableWriter) formatString() string {
	f := []string{}

	ws := t.widths()

	for _, w := range ws {
		f = append(f, fmt.Sprintf("%%-%ds", w))
	}

	return strings.Join(f, "  ")
}

func (t *tableWriter) printJSON() error {
	lccs := []string{}

	for _, c := range t.columns {
		lccs = append(lccs, strings.ToLower(c.(string)))
	}

	vs := []map[string]any{}

	for _, r := range t.rows {
		v := map[string]any{}

		for i := range t.columns {
			v[lccs[i]] = r[i]
		}

		vs = append(vs, v)
	}

	data, err := json.MarshalIndent(vs, "", "  ")
	if err != nil {
		return err
	}

	if _, err := t.ctx.Write(data); err != nil {
		return err
	}

	return nil
}

func (t *tableWriter) printText() error {
	f := t.formatString()

	t.ctx.Writef(fmt.Sprintf("<h1>%s</h1>\n", f), t.columns...) //nolint:errcheck

	for _, r := range t.rows {
		t.ctx.Writef(fmt.Sprintf("<value>%s</value>\n", f), r...) //nolint:errcheck
	}

	return nil
}

func (t *tableWriter) widths() []int {
	w := make([]int, len(t.columns))

	for i, c := range t.columns {
		w[i] = len(stripTag(c))

		for _, r := range t.rows {
			if len(r) > i {
				if lri := len(stripTag(r[i])); lri > w[i] {
					w[i] = lri
				}
			}
		}
	}

	w[len(w)-1] = 0

	return w
}
