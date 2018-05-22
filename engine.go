package stdcli

import (
	"strings"

	"github.com/spf13/pflag"
)

type Engine struct {
	Commands []Command
	Name     string
	Version  string
	Writer   *Writer
}

func (e *Engine) Command(command, description string, fn HandlerFunc, opts CommandOptions) {
	e.Commands = append(e.Commands, Command{
		Command:     strings.Split(command, " "),
		Description: description,
		Handler:     fn,
		Flags:       opts.Flags,
		Validate:    opts.Validate,
		engine:      e,
	})
}

func (e *Engine) Execute(args []string) int {
	fs := pflag.NewFlagSet(e.Name, pflag.ContinueOnError)
	fs.Usage = func() {}
	fs.Parse(args)

	var m *Command
	var cargs []string

	for _, c := range e.Commands {
		d := c
		if a, ok := d.Match(args); ok {
			if m == nil || len(m.Command) < len(c.Command) {
				m = &d
				cargs = a
			}
		}
	}

	if m == nil {
		m = &(e.Commands[0])
	}

	err := m.Execute(cargs)
	switch t := err.(type) {
	case nil:
		return 0
	default:
		e.Writer.Errorf("%s", t)
		return 1
	}

	return 0
}
