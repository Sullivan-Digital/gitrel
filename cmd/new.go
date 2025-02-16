package cmd

import (
	"errors"
	"gitrel/git"
	"gitrel/semver"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new release branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCtx := git.NewCmdGitContext()
		ctx, err := getCommandContext(gitCtx)
		if err != nil {
			return err
		}

		if len(args) != 1 {
			return cmd.Help()
		}

		version := args[0]
		if !semver.ValidateSemver(version) {
			return errors.New("invalid version format. please use semantic versioning (e.g., 1.0.0, 1.2.3-alpha, 2.0.0+build.1)")
		}

		git.CreateReleaseBranch(version, ctx, gitCtx)

		return nil
	},
}

func init() {
	newCmd.AddCommand(newVersionCmd)
	newCmd.AddCommand(newMajorCmd)
	newCmd.AddCommand(newMinorCmd)
	newCmd.AddCommand(newPatchCmd)
}
