package command

import (
	"github.com/spf13/cobra"
)

type Command struct {
	cmd       *cobra.Command
	initFuncs []func() error
}

func New(cmd *cobra.Command, initFuncs ...func() error) *Command {
	return &Command{cmd: cmd, initFuncs: initFuncs}
}

func (c Command) Init() {
	for _, initFunc := range c.initFuncs {
		cobra.CheckErr(initFunc())
	}
}

func (c Command) Cmd() *cobra.Command {
	return c.cmd
}
