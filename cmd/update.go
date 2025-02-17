package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Push changes to a release branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		runUpdateCmd(args, ctx)
		return nil
	},
}

var updateVersionCmd = &cobra.Command{
	Use:   "<version>",
	Short: "Push changes to the release branch matching the specified version prefix",
}

func init() {
	updateCmd.AddCommand(updateVersionCmd)
}

func runUpdateCmd(args []string, ctx interfaces.GitRelContext) {
	err := git.UpdateVersion(args[0], ctx)
	if err != nil {
		ctx.Output().Println(err)
	}
}
