package stdcli

import "os/exec"

type Executor interface {
	Run(cmd string, args ...string) ([]byte, error)
}

type CmdExecutor struct {
}

func (e *CmdExecutor) Run(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}
