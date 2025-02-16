package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"
	"gitrel/semver"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new release branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		runNewCmd(cmd, args, ctx)
		return nil
	},
}

func init() {
	newCmd.AddCommand(newVersionCmd)
	newCmd.AddCommand(newMajorCmd)
	newCmd.AddCommand(newMinorCmd)
	newCmd.AddCommand(newPatchCmd)
}

func runNewCmd(cmd *cobra.Command, args []string, ctx interfaces.GitRelContext) {
	if len(args) != 1 {
		cmd.Help()
	}

	version := args[0]
	if !semver.ValidateSemver(version) {
		ctx.Output().Println("invalid version format. please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)")
	}

	git.CreateReleaseBranch(version, ctx)
}
