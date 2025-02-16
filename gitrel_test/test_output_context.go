package gitrel_test

import "fmt"

type TestOutputContext struct {
	Output string
}

func DefaultTestOutputContext() *TestOutputContext {
	return &TestOutputContext{
		Output: "",
	}
}

func (c *TestOutputContext) Print(args ...interface{}) {
	c.Output += fmt.Sprint(args...)
}

func (c *TestOutputContext) Println(args ...interface{}) {
	c.Output += fmt.Sprintln(args...)
}

func (c *TestOutputContext) Printf(format string, args ...interface{}) {
	c.Output += fmt.Sprintf(format, args...)
}





