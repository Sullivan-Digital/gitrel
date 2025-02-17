package cmd

import (
	"gitrel/git"
	"gitrel/interfaces"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [<version> | latest]",
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

var updateLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Push changes to the latest release branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := NewCmdGitRelContext()
		if err != nil {
			return err
		}

		runUpdateCmd([]string{"latest"}, ctx)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateVersionCmd)
	updateCmd.AddCommand(updateLatestCmd)
}

func runUpdateCmd(args []string, ctx interfaces.GitRelContext) {
	err := git.UpdateVersion(args[0], ctx)
	if err != nil {
		ctx.Output().Println(err)
	}
}
