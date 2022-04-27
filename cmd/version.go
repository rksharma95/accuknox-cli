package cmd

import (
	"github.com/accuknox/accuknox-cli/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the get command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display version information`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//Print KubeArmor version information
		if err := version.PrintVersion(client); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
