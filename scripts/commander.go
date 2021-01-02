package scripts

import "os/exec"

type Commander interface {
	Command(name string, arg ...string) *exec.Cmd
}

type StandardCommander struct{}

func (sc *StandardCommander) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

var _ Commander = (*StandardCommander)(nil)
