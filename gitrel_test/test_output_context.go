package gitrel_test

import (
	"fmt"
	"strings"
	"testing"
)

type TestOutputContext struct {
	Output  string
	testCtx *testing.T
}

func DefaultTestOutputContext(t *testing.T) *TestOutputContext {
	return &TestOutputContext{
		Output:  "",
		testCtx: t,
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

func (c *TestOutputContext) AssertOutput(expectedOutput string) {
	c.testCtx.Helper()
	if c.Output != expectedOutput {
		printMismatchedOutput(c.testCtx, expectedOutput, c.Output)
	}
}

func (c *TestOutputContext) AssertOutputLines(expectedOutputLines ...string) {
	c.testCtx.Helper()
	c.AssertOutput(strings.TrimSpace(strings.Join(expectedOutputLines, "\n")) + "\n")
}

func printMismatchedOutput(t *testing.T, expectedOutput string, actualOutput string) {
	t.Helper()
	t.Fatalf(strings.Join([]string{
		"Inconsistent output - expected vs. actual:",
		"----------------expected----------------",
		normaliseOutputForPrinting(expectedOutput),
		"-----------------actual-----------------",
		normaliseOutputForPrinting(actualOutput),
		"----------------------------------------",
	}, "\n"))
}

func normaliseOutputForPrinting(output string) string {
	return strings.TrimSpace(strings.ReplaceAll(output, "\n", "<br>\n"))
}
