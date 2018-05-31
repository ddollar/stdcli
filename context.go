package stdcli

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

type Context struct {
	Args  []string
	Flags []*Flag

	engine *Engine
	state  *terminal.State
}

func (c *Context) Arg(i int) string {
	if i < len(c.Args) {
		return c.Args[i]
	}

	return ""
}

func (c *Context) Flag(name string) *Flag {
	for _, f := range c.Flags {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (c *Context) Bool(name string) bool {
	if f := c.Flag(name); f != nil && f.Kind == reflect.Bool {
		switch t := f.Value.(type) {
		case nil:
			v, _ := f.Default.(bool)
			return v
		case bool:
			return t
		}
	}
	return false
}

func (c *Context) Int(name string) int {
	if f := c.Flag(name); f != nil && f.Kind == reflect.Int {
		switch t := f.Value.(type) {
		case nil:
			v, _ := f.Default.(int)
			return v
		case int:
			return t
		}
	}
	return 0
}

func (c *Context) String(name string) string {
	if f := c.Flag(name); f != nil && f.Kind == reflect.String {
		switch t := f.Value.(type) {
		case nil:
			v, _ := f.Default.(string)
			return v
		case string:
			return t
		}
	}
	return ""
}

func (c *Context) Info() *Info {
	return &Info{Context: c}
}

func (c *Context) ModeRaw() error {
	state, err := terminal.MakeRaw(int(c.Reader().Fd()))
	if err != nil {
		return err
	}

	c.state = state

	return nil
}

func (c *Context) ModeRestore() error {
	if c.state == nil {
		return nil
	}

	return terminal.Restore(int(c.Reader().Fd()), c.state)
}

func (c *Context) Reader() *Reader {
	return c.engine.Reader
}

func (c *Context) Table(columns ...string) *Table {
	return &Table{Columns: columns, Context: c}
}

func (c *Context) Write(data []byte) (int, error) {
	return c.Writer().Write(data)
}

func (c *Context) Writer() *Writer {
	return c.engine.Writer
}

func (c *Context) OK(id ...string) error {
	c.Writer().Writef("<ok>OK</ok>")

	if len(id) > 0 {
		c.Writer().Writef(", <id>%s</id>", strings.Join(id, " "))
	}

	c.Writer().Writef("\n")

	return nil
}

func (c *Context) Startf(format string, args ...interface{}) {
	c.Writer().Writef(fmt.Sprintf("%s... ", format), args...)
}

func (c *Context) Writef(format string, args ...interface{}) error {
	_, err := c.Writer().Writef(format, args...)
	return err
}

func (c *Context) Options(opts interface{}) error {
	v := reflect.ValueOf(opts).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		u := v.Field(i)

		if n := f.Tag.Get("flag"); n != "" {
			switch f.Type.Elem().Kind() {
			case reflect.Bool:
				var x bool
				y := c.Bool(strings.Split(n, ",")[0])
				if x != y {
					u.Set(reflect.ValueOf(&y))
				}
			case reflect.Int:
				var x int
				y := c.Int(strings.Split(n, ",")[0])
				if x != y {
					u.Set(reflect.ValueOf(&y))
				}
			case reflect.String:
				var x string
				y := c.String(strings.Split(n, ",")[0])
				if x != y {
					u.Set(reflect.ValueOf(&y))
				}
			default:
				return fmt.Errorf("unknown flag type: %s", f.Type.Elem().Kind())
			}
		}
	}

	return nil
}
