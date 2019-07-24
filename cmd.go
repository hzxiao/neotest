package neotest

import (
	"github.com/spf13/cobra"
)

type Commander interface {
	Exec(vm *VM) error
	Line() int
	Name() string
}

type Cmd struct {
	name string
	line int
	cmd  *cobra.Command
}

func NewCmd(name string, desc string, line int) *Cmd {
	return &Cmd{
		name: name,
		line:line,
		cmd:  &cobra.Command{Use: name, Short: desc},
	}
}

func (c *Cmd) Run(vm *VM) error {
	panic("must impl")
}

func (c *Cmd) Line() int {
	return c.line
}

func (c *Cmd) Name() string {
	return c.name
}