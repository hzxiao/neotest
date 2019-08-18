package main

import (
	"fmt"
	"github.com/hzxiao/neotest"
	"github.com/hzxiao/neotest/pkg/pln"
	"github.com/spf13/cobra"
)

var verbose bool

func main() {
	root := &cobra.Command{
		Use:   "neotest",
		Short: "An auto tool for testing neo transaction",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				return
			}
			err := run(args)
			if err != nil {
				pln.Error(err)
			}
		},
	}
	root.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")

	root.Execute()
}

func run(files []string) error {
	pln.Verbose = verbose

	for _, file := range files {
		src, err := neotest.NewSource(file)
		if err != nil {
			return err
		}
		commands, err := src.Parse()
		if err != nil {
			return fmt.Errorf("parse %v err: %v", file, err)
		}
		vm := neotest.NewVM(commands)
		err = vm.Run()
		if err != nil {
			return fmt.Errorf("run %v err: %v", file, err)
		}
	}
	return nil
}
