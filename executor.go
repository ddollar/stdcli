package stdcli

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type Executor interface {
	Execute(cmd string, args ...string) ([]byte, error)
	Run(w io.Writer, cmd string, args ...string) error
	Terminal(cmd string, args ...string) error
}

type CmdExecutor struct {
}

func (e *CmdExecutor) Execute(cmd string, args ...string) ([]byte, error) {
	data, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func (e *CmdExecutor) Run(w io.Writer, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)

	c.Stdout = w
	c.Stderr = w

	if err := c.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (e *CmdExecutor) Terminal(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
