package stdcli

import (
	"fmt"
)

type Info struct {
	Context *Context
	Rows    []InfoRow
}

type InfoRow struct {
	Header string
	Value  string
}

func (i *Info) Add(header, value string) {
	i.Rows = append(i.Rows, InfoRow{Header: header, Value: value})
}

func (i *Info) Print() error {
	f := i.formatString()

	for _, r := range i.Rows {
		i.Context.Writef(f, r.Header, r.Value)
	}

	return nil
}

func (i *Info) formatString() string {
	return fmt.Sprintf("<h1>%%-%ds</h1>  <value>%%-%ds</value>\n", i.headerWidth(), i.valueWidth())
}

func (i *Info) headerWidth() int {
	w := 0

	for _, r := range i.Rows {
		if len(r.Header) > w {
			w = len(r.Header)
		}
	}

	return w
}

func (i *Info) valueWidth() int {
	w := 0

	for _, r := range i.Rows {
		if len(r.Value) > w {
			w = len(r.Value)
		}
	}

	return w
}
