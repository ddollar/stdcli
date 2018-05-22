package stdcli

import (
	"fmt"
)

type Validator func(c *Context) error

func Args(num int) Validator {
	return func(c *Context) error {
		if len(c.Args) != num {
			return fmt.Errorf("%d args required", num)
		}
		return nil
	}
}

func ArgsBetween(min, max int) Validator {
	return func(c *Context) error {
		if len(c.Args) < min {
			return fmt.Errorf("at least %d args required", min)
		}
		if len(c.Args) > max {
			return fmt.Errorf("too many args")
		}
		return nil
	}
}
