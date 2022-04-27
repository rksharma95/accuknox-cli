package cmd

import (
	"fmt"

	ciliumk8s "github.com/cilium/cilium-cli/k8s"
	"github.com/kubearmor/kubearmor-client/k8s"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	contextName string

	client    *k8s.Client
	k8sClient *ciliumk8s.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		//Initialise k8sClient for all child commands to inherit
		client, err = k8s.ConnectK8sClient()
		// fmt.Printf("%v", client.K8sClientset)
		if err != nil {
			log.Error().Msgf("unable to create Kubernetes clients: %s", err.Error())
			return err
		}

		c, err := ciliumk8s.NewClient(contextName, "")
		if err != nil {
			return fmt.Errorf("unable to create Kubernetes client: %w", err)
		}

		k8sClient = c

		return nil
	},
	Use:   "accuknox",
	Short: "CLI Utility to help manage Accuknox security solution",
	Long: `CLI Utility to help manage Accuknox security solution
	accuknox cli tool helps to install, manage and troubleshoot Accuknox security solution
	`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
