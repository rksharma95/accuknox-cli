package cmd

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// portForwardCmd represents the accuknox port-forward command
var portForwardCmd = &cobra.Command{
	Use:   "port-forward",
	Short: "port-forward KubeArmor, Cilium and Discovery-engine in a Kubernetes Cluster",
	Long:  `port-forward KubeArmor, Cilium and Discovery-engine in a Kubernetes Clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}

		return nil
	},
}

// accuknox port-forward kubearmor
var karmorCmd = &cobra.Command{
	Use:   "kubearmor",
	Short: "port-forward KubeArmor in a Kubernetes Cluster",
	Long:  `port-forward KubeArmor relay port to local machine in a Kubernetes Clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {

		cmdArgs := []string{
			"-n",
			"kube-system",
			"port-forward",
			"service/kubearmor",
			"--address",
			"0.0.0.0",
			"--address",
			"::",
			"32767:32767"}
		pfCmd := exec.Command("kubectl", cmdArgs...)

		bytes, err := pfCmd.CombinedOutput()

		if err != nil {
			return errors.New("Unable to port-forward kubearmor: " + err.Error())
		}

		fmt.Println(string(bytes))
		return nil
	},
}

// accuknox port-forward cilium
var ciliumCmd = &cobra.Command{
	Use:   "cilium",
	Short: "port-forward Cilium in a Kubernetes Cluster",
	Long:  `port-forward Cilium hubble-relay port to local machine in a Kubernetes Clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {

		cmdArgs := []string{
			"-n",
			"kube-system",
			"port-forward",
			"service/hubble-relay",
			"--address",
			"0.0.0.0",
			"--address",
			"::",
			"4245:80"}

		pfCmd := exec.Command("kubectl", cmdArgs...)

		bytes, err := pfCmd.CombinedOutput()

		if err != nil {
			return errors.New("Unable to port-forward cilium: " + err.Error())
		}

		fmt.Println(string(bytes))

		return nil
	},
}

// accuknox port-forward discovery-engine
var dEngineCmd = &cobra.Command{
	Use:   "discovery-engine",
	Short: "port-forward Discovery-engine in a Kubernetes Cluster",
	Long:  `port-forward Discovery-engine in a Kubernetes Clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {

		cmdArgs := []string{
			"-n",
			"explorer",
			"port-forward",
			"service/knoxautopolicy",
			"--address",
			"0.0.0.0",
			"--address",
			"::",
			"9089:9089"}

		pfCmd := exec.Command("kubectl", cmdArgs...)

		bytes, err := pfCmd.CombinedOutput()

		if err != nil {
			return errors.New("Unable to port-forward discovery engine: " + err.Error())
		}

		fmt.Println(string(bytes))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(portForwardCmd)

	portForwardCmd.AddCommand(karmorCmd)
	portForwardCmd.AddCommand(ciliumCmd)
	portForwardCmd.AddCommand(dEngineCmd)

}
