package cmd

import (
	"context"
	"os"
	"time"

	di "github.com/accuknox/accuknox-cli/install"
	"github.com/cilium/cilium-cli/defaults"
	ci "github.com/cilium/cilium-cli/install"
	ki "github.com/kubearmor/kubearmor-client/install"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	namespace      string
	installOptions ki.Options
	params         = ci.Parameters{Writer: os.Stdout}
	diOptions      di.Options
)

const (
	encryptionDisabled  = "disabled"
	encryptionIPsec     = "ipsec"
	encryptionWireguard = "wireguard"
)

// installCmd represents the get command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install KubeArmor, Cilium and Discovery-engine in a Kubernetes Cluster",
	Long:  `Install KubeArmor, Cilium and Discovery-engine in a Kubernetes Clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Install Cilium
		params.Namespace = namespace
		installer, err := ci.NewK8sInstaller(k8sClient, params)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		if err := installer.Install(context.Background()); err != nil {
			installer.RollbackInstallation(context.Background())

			log.Error().Msgf("Unable to install Cilium: %s", err.Error())
		}

		// Install KubeArmor
		installOptions.Namespace = namespace
		if err := ki.K8sInstaller(client, installOptions); err != nil {
			return err
		}

		// Install dscovery-engine
		diOptions.Namespace = "explorer"
		if err := di.DiscoveryEngineInstaller(client, diOptions); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&installOptions.KubearmorImage, "image", "i", "kubearmor/kubearmor:stable", "Kubearmor daemonset image to use")
	installCmd.Flags().StringVarP(&namespace, "namespace", "n", "kube-system", "Namespace for resources")
	installCmd.Flags().StringVar(&params.ClusterName, "cluster-name", "", "Name of the cluster")
	installCmd.Flags().StringSliceVar(&params.DisableChecks, "disable-check", []string{}, "Disable a particular validation check")
	installCmd.Flags().StringVar(&params.Version, "version", defaults.Version, "Cilium version to install")
	installCmd.Flags().StringVar(&params.BaseVersion, "base-version", defaults.Version,
		"Specify the base Cilium version for configuration purpose in case the --version flag doesn't indicate the actual Cilium version")
	installCmd.Flags().MarkHidden("base-version")
	installCmd.Flags().StringVar(&params.DatapathMode, "datapath-mode", "", "Datapath mode to use")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.IPAM, "ipam", "", "IP Address Management (IPAM) mode")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.IPv4NativeRoutingCIDR, "ipv4-native-routing-cidr", "", "IPv4 CIDR within which native routing is possible")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().IntVar(&params.ClusterID, "cluster-id", 0, "Unique cluster identifier for multi-cluster")
	installCmd.Flags().StringVar(&params.InheritCA, "inherit-ca", "", "Inherit/import CA from another cluster")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.KubeProxyReplacement, "kube-proxy-replacement", "disabled", "Enable/disable kube-proxy replacement { disabled | probe | strict }")
	installCmd.Flags().BoolVar(&params.Wait, "wait", true, "Wait for status to report success (no errors)")
	installCmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", defaults.StatusWaitDuration, "Maximum time to wait for status")
	installCmd.Flags().BoolVar(&params.RestartUnmanagedPods, "restart-unmanaged-pods", true, "Restart pods which are not being managed by Cilium")
	installCmd.Flags().StringVar(&params.Encryption, "encryption", "disabled", "Enable encryption of all workloads traffic { disabled | ipsec | wireguard }")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().BoolVar(&params.NodeEncryption, "node-encryption", false, "Enable encryption of all node to node traffic")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringSliceVar(&params.ConfigOverwrites, "config", []string{}, "Set ConfigMap entries { key=value[,key=value,..] }")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.AgentImage, "agent-image", "", "Image path to use for Cilium agent")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.OperatorImage, "operator-image", "", "Image path to use for Cilium operator")
	installCmd.Flags().DurationVar(&params.CiliumReadyTimeout, "cilium-ready-timeout", 5*time.Minute,
		"Timeout for Cilium to become ready before restarting unmanaged pods")
	installCmd.Flags().BoolVar(&params.Rollback, "rollback", true, "Roll back installed resources on failure")

	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.Azure.ResourceGroupName, "azure-resource-group", "", "Azure resource group name the cluster is in (required)")
	installCmd.Flags().StringVar(&params.Azure.AKSNodeResourceGroup, "azure-node-resource-group", "", "Azure node resource group name the cluster is in. Bypasses `--azure-resource-group` if provided.")
	installCmd.Flags().MarkHidden("azure-node-resource-group") // intended for for development purposes, notably CI usage, cf. azure.go
	installCmd.Flags().StringVar(&params.Azure.SubscriptionName, "azure-subscription", "", "Azure subscription name the cluster is in (default `az account show`)")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.Azure.SubscriptionID, "azure-subscription-id", "", "Azure subscription ID. Bypasses auto-detection and `--azure-subscription` if provided.")
	installCmd.Flags().MarkHidden("azure-subscription-id") // intended for for development purposes, notably CI usage, cf. azure.go
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.Azure.TenantID, "azure-tenant-id", "", "Tenant ID of Azure Service Principal to use for installing Cilium (will create one if none provided)")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.Azure.ClientID, "azure-client-id", "", "Client (application) ID of Azure Service Principal to use for installing Cilium (will create one if none provided)")
	// It can be deprecated since we have a helm option for it
	installCmd.Flags().StringVar(&params.Azure.ClientSecret, "azure-client-secret", "", "Client secret of Azure Service Principal to use for installing Cilium (will create one if none provided)")
	installCmd.Flags().StringVar(&params.K8sVersion, "k8s-version", "", "Kubernetes server version in case auto-detection fails")

	installCmd.Flags().StringVar(&params.HelmChartDirectory, "chart-directory", "", "Helm chart directory")
	installCmd.Flags().StringSliceVar(&params.HelmOpts.ValueFiles, "helm-values", []string{}, "Specify helm values in a YAML file or a URL (can specify multiple)")
	installCmd.Flags().StringArrayVar(&params.HelmOpts.Values, "helm-set", []string{}, "Set helm values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	installCmd.Flags().StringArrayVar(&params.HelmOpts.StringValues, "helm-set-string", []string{}, "Set helm STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	installCmd.Flags().StringArrayVar(&params.HelmOpts.FileValues, "helm-set-file", []string{}, "Set helm values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)")
	installCmd.Flags().StringVar(&params.HelmGenValuesFile, "helm-auto-gen-values", "", "Write an auto-generated helm values into this file")
	installCmd.Flags().StringVar(&params.HelmValuesSecretName, "helm-values-secret-name", defaults.HelmValuesSecretName, "Secret name to store the auto-generated helm values file. The namespace is the same as where Cilium will be installed")
	installCmd.Flags().StringVar(&params.ImageSuffix, "image-suffix", "", "Set all generated images with this suffix")
	installCmd.Flags().StringVar(&params.ImageTag, "image-tag", "", "Set all images with this tag")

}
