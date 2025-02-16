package cmd

import (
	"github.com/spf13/cobra"
)

var newVersionCmd = &cobra.Command{
	Use:   "<version>",
	Short: "Create a new release branch with the specified version",
}