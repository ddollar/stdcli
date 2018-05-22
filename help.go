package stdcli

import (
	"fmt"
	"sort"
)

func Help(c *Context) error {
	helpGlobal(c.engine)
	return nil
}

func helpGlobal(e *Engine) {
	cs := []Command{}

	for _, cmd := range e.Commands {
		cs = append(cs, cmd)
	}

	sort.Slice(cs, func(i, j int) bool { return cs[i].FullCommand() < cs[j].FullCommand() })

	l := 0

	for _, cmd := range cs {
		c := cmd.FullCommand()

		if len(c) > l {
			l = len(c)
		}
	}

	for _, cmd := range cs {
		e.Writer.Writef(fmt.Sprintf(fmt.Sprintf("<h1>%%-%ds</h1>  <value>%%s</value>\n", l), cmd.FullCommand(), cmd.Description))
	}
}

func helpCommand(e *Engine, cmd *Command) {
	e.Writer.Writef("<h1>%s</h1>\n\n", cmd.FullCommand())

	e.Writer.Writef("<h2>description</h2>\n  <value>%s</value>\n\n", cmd.Description)

	e.Writer.Writef("<h2>options</h2>\n")

	ll := 0
	ls := 0

	for _, f := range cmd.Flags {
		l := f.UsageLong()
		s := f.UsageShort()

		if len(l) > ll {
			ll = len(l)
		}

		if len(s) > ls {
			ls = len(s)
		}
	}

	for _, f := range cmd.Flags {
		// e.Writer.Writef(fmt.Sprintf(fmt.Sprintf("  %%-%ds  %%-%ds   <header>%%s</h1>\n", ll, ls), f.UsageLong(), f.UsageShort(), f.Description))
		e.Writer.Writef(fmt.Sprintf(fmt.Sprintf("  %%-%ds  %%-%ds\n", ll, ls), f.UsageLong(), f.UsageShort()))
	}
}
