package stdcli

import (
	"fmt"
	"reflect"
	"strings"
)

type Context struct {
	Args  []string
	Flags []*Flag

	engine *Engine
}

func (c *Context) Arg(i int) string {
	if i < len(c.Args) {
		return c.Args[i]
	}

	return ""
}

func (c *Context) String(name string) string {
	for _, f := range c.Flags {
		if f.Name == name && f.Kind == reflect.String {
			switch t := f.Value.(type) {
			case nil:
				v, _ := f.Default.(string)
				return v
			case string:
				return t
			default:
				return ""
			}
		}
	}

	return ""
}

func (c *Context) Info() *Info {
	return &Info{Context: c}
}

func (c *Context) Table(columns ...string) *Table {
	return &Table{Columns: columns, Context: c}
}

func (c *Context) Writer() *Writer {
	return c.engine.Writer
}

func (c *Context) OK() error {
	c.Writer().Writef("<ok>OK</ok>\n")
	return nil
}

func (c *Context) Startf(format string, args ...interface{}) {
	c.Writer().Writef(fmt.Sprintf("%s... ", format), args...)
}

func (c *Context) Writef(format string, args ...interface{}) {
	c.Writer().Writef(format, args...)
}

func (c *Context) Options(opts interface{}) error {
	v := reflect.ValueOf(opts).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		u := v.Field(i)

		if n := f.Tag.Get("flag"); n != "" {
			switch f.Type.Elem().Kind() {
			case reflect.String:
				s := c.String(strings.Split(n, ",")[0])
				u.Set(reflect.ValueOf(&s))
			default:
				return fmt.Errorf("unknown flag type: %s", f.Type.Elem().Kind())
			}
		}
	}

	return nil
}
