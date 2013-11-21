package gist7576154

import (
	"os/exec"
)

type CmdTemplate struct {
	name string
	args []string
	Dir  string
}

func NewCmdTemplate(name string, arg ...string) CmdTemplate {
	return CmdTemplate{
		name: name,
		args: arg,
	}
}

func (ct CmdTemplate) NewCommand() *exec.Cmd {
	cmd := exec.Command(ct.name, ct.args...)
	cmd.Dir = ct.Dir
	return cmd
}
