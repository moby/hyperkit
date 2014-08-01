package gist7576154

import (
	"io"
	"os/exec"

	. "github.com/shurcooL/go/gists/gist7729255"
	. "github.com/shurcooL/go/gists/gist7802150"

	"gopkg.in/pipe.v2"
)

// CmdFactory is an interface to create new commands.
type CmdFactory interface {
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

// TODO: See if there's some way to initialize DepNode2Func.UpdateFunc through NewCmdTemplateDynamic2().
func NewCmdTemplateDynamic2() *CmdTemplateDynamic2 {
	return &CmdTemplateDynamic2{}
}

func (this *CmdTemplateDynamic2) NewCommand() *exec.Cmd {
	MakeUpdated(this)
	return this.Template.NewCommand()
}

// =====

type PipeFactory interface {
	NewPipe(stdout, stderr io.Writer) (*pipe.State, pipe.Pipe)
}

// ---

type PipeStatic pipe.Pipe

func (this PipeStatic) NewPipe(stdout, stderr io.Writer) (*pipe.State, pipe.Pipe) {
	return pipe.NewState(stdout, stderr), (pipe.Pipe)(this)
}

// ---

type pipeTemplate struct {
	Pipe  pipe.Pipe
	Dir   string
	Stdin func() io.Reader
}

func NewPipeTemplate(pipe pipe.Pipe) *pipeTemplate {
	return &pipeTemplate{Pipe: pipe}
}

func (this *pipeTemplate) NewPipe(stdout, stderr io.Writer) (*pipe.State, pipe.Pipe) {
	s := pipe.NewState(stdout, stderr)
	s.Dir = this.Dir
	if this.Stdin != nil {
		s.Stdin = this.Stdin()
	}
	return s, this.Pipe
}

// ---

type pipeTemplateDynamic struct {
	Template *pipeTemplate

	DepNode2Func
}

func NewPipeTemplateDynamic() *pipeTemplateDynamic {
	return &pipeTemplateDynamic{}
}

func (this *pipeTemplateDynamic) NewPipe(stdout, stderr io.Writer) (*pipe.State, pipe.Pipe) {
	MakeUpdated(this)
	return this.Template.NewPipe(stdout, stderr)
}
