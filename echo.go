package neotest

import (
	"bufio"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
)

type EchoCmd struct {
	*Cmd
}

func NewEchoCmd(line int) *EchoCmd {
	return &EchoCmd{
		Cmd: NewCmd("echo", "echo <object1> <object2>", line),
	}
}

func parseEchoCmd(line int, buf *bufio.Scanner) (*EchoCmd, error) {
	return nil, nil
}

func (echo *EchoCmd) Exec(vm *VM) error  {
	echo.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			pp.Printf("%v ", arg)
		}
		pp.Println()
		return nil
	}
	return nil
}