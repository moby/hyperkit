package gist7576154

import (
	"io"
	"os/exec"

	. "gist.github.com/7729255.git"
	. "gist.github.com/7802150.git"
)

// CmdTemplater is the interface for command templates.
type CmdTemplater interface {
	NewCommand() *exec.Cmd
}

// CmdTemplate is a command template.
type CmdTemplate struct {
	NameArgs []string
	Dir      string
	Stdin    func() io.Reader
}

// NewCmdTemplate returns a CmdTemplate.
func NewCmdTemplate(name string, arg ...string) CmdTemplate {
	return CmdTemplate{
		NameArgs: append([]string{name}, arg...),
	}
}

// NewCommand generates a new *exec.Cmd from the template.
func (ct CmdTemplate) NewCommand() *exec.Cmd {
	cmd := exec.Command(ct.NameArgs[0], ct.NameArgs[1:]...)
	cmd.Dir = ct.Dir
	if ct.Stdin != nil {
		cmd.Stdin = ct.Stdin()
	}
	return cmd
}

// ---

type CmdTemplateDynamic struct {
	NameArgs Strings
	Dir      String
	Stdin    func() io.Reader
}

func NewCmdTemplateDynamic(nameArgs Strings) CmdTemplateDynamic {
	return CmdTemplateDynamic{
		NameArgs: nameArgs,
	}
}

func (ct CmdTemplateDynamic) NewCommand() *exec.Cmd {
	nameArgs := ct.NameArgs.Get()
	cmd := exec.Command(nameArgs[0], nameArgs[1:]...)
	if ct.Dir != nil {
		cmd.Dir = ct.Dir.Get()
	}
	if ct.Stdin != nil {
		cmd.Stdin = ct.Stdin()
	}
	return cmd
}

// ---

type CmdTemplateDynamic2 struct {
	Template CmdTemplate

	DepNode2Func
}

func NewCmdTemplateDynamic2() *CmdTemplateDynamic2 {
	return &CmdTemplateDynamic2{}
}

func (this *CmdTemplateDynamic2) NewCommand() *exec.Cmd {
	MakeUpdated(this)
	return this.Template.NewCommand()
}
