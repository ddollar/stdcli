package stdcli

import (
	"fmt"
	"sort"
)

func help(e *Engine) HandlerFunc {
	return func(ctx Context) error {
		cs := []Command{}

		for _, cmd := range e.Commands {
			if cmd.Invisible {
				continue
			}

			cs = append(cs, cmd)
		}

		sort.Slice(cs, func(i, j int) bool { return cs[i].FullCommand() < cs[j].FullCommand() })

		l := 7

		for _, cmd := range cs {
			c := cmd.FullCommand()

			if len(c) > l {
				l = len(c)
			}
		}

		for _, cmd := range cs {
			e.Writer.Writef(fmt.Sprintf(fmt.Sprintf("<h1>%%-%ds</h1>  <value>%%s</value>\n", l), cmd.FullCommand(), cmd.Description)) // nolint:errcheck
		}

		return nil
	}
}

func helpCommand(ctx Context, e *Engine, cmd *Command) {
	e.Writer.Writef("<h2>USAGE</h2>\n  <value>%s</value> <info>%s</info>\n\n", cmd.FullCommand(), cmd.Usage) //nolint:errcheck
	e.Writer.Writef("<h2>DESCRIPTION</h2>\n  <value>%s</value>\n\n", cmd.Description)                        //nolint:errcheck
	e.Writer.Writef("<h2>OPTIONS</h2>\n")                                                                    //nolint:errcheck

	fs := []Flag(cmd.Flags)

	cw := ctx.Columns()

	for _, f := range fs {
		cw.Append("", f.Usage(), f.Description)
	}

	cw.Print()

	// for _, f := range fs {
	// 	// e.Writer.Writef(fmt.Sprintf(fmt.Sprintf("  %%-%ds  %%s\n", maxusage), f.Usage(), f.Description)) //nolint:errcheck
	// 	e.Writer.WriteColumns(f.Usage(), f.Description)
	// }
}
