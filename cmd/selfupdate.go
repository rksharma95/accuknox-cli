package cmd

import (
	"github.com/accuknox/accuknox-cli/selfupdate"
	"github.com/spf13/cobra"
)

// selfUpdateCmd represents the get command
var selfUpdateCmd = &cobra.Command{
	Use:   "selfupdate",
	Short: "selfupdate this cli tool",
	Long:  `selfupdate this cli tool for checking the latest release on the github`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//Print KubeArmor version information
		if err := selfupdate.SelfUpdate(client); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}
