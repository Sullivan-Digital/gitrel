package cmd

import (
	"fmt"
	"gitrel/interfaces"
)

type CmdOutputContext struct{}

func NewCmdOutputContext() interfaces.OutputContext {
	return &CmdOutputContext{}
}

func (c *CmdOutputContext) Print(args ...interface{}) {
	fmt.Print(args...)
}

func (c *CmdOutputContext) Println(args ...interface{}) {
	fmt.Println(args...)
}

func (c *CmdOutputContext) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
