package cmd

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/cilium/cilium-cli/defaults"
	ci "github.com/cilium/cilium-cli/install"
	"github.com/cilium/cilium-cli/k8s"
	"github.com/cilium/cilium/api/v1/models"
	ki "github.com/kubearmor/kubearmor-client/install"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/blang/semver/v4"
	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

var (
	namespace      string
	installOptions ki.Options
	params         = ci.Parameters{Writer: os.Stdout}
)

const (
	encryptionDisabled  = "disabled"
	encryptionIPsec     = "ipsec"
	encryptionWireguard = "wireguard"
)

type k8sInstallerImplementation interface {
	ClusterName() string
	ListNodes(ctx context.Context, options metav1.ListOptions) (*corev1.NodeList, error)
	GetCiliumExternalWorkload(ctx context.Context, name string, opts metav1.GetOptions) (*ciliumv2.CiliumExternalWorkload, error)
	CreateCiliumExternalWorkload(ctx context.Context, cew *ciliumv2.CiliumExternalWorkload, opts metav1.CreateOptions) (*ciliumv2.CiliumExternalWorkload, error)
	DeleteCiliumExternalWorkload(ctx context.Context, name string, opts metav1.DeleteOptions) error
	ListCiliumExternalWorkloads(ctx context.Context, opts metav1.ListOptions) (*ciliumv2.CiliumExternalWorkloadList, error)
	CreateServiceAccount(ctx context.Context, namespace string, account *corev1.ServiceAccount, opts metav1.CreateOptions) (*corev1.ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error
	GetConfigMap(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ConfigMap, error)
	CreateConfigMap(ctx context.Context, namespace string, config *corev1.ConfigMap, opts metav1.CreateOptions) (*corev1.ConfigMap, error)
	DeleteConfigMap(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error
	CreateClusterRole(ctx context.Context, config *rbacv1.ClusterRole, opts metav1.CreateOptions) (*rbacv1.ClusterRole, error)
	DeleteClusterRole(ctx context.Context, name string, opts metav1.DeleteOptions) error
	CreateClusterRoleBinding(ctx context.Context, role *rbacv1.ClusterRoleBinding, opts metav1.CreateOptions) (*rbacv1.ClusterRoleBinding, error)
	DeleteClusterRoleBinding(ctx context.Context, name string, opts metav1.DeleteOptions) error
	CreateDaemonSet(ctx context.Context, namespace string, ds *appsv1.DaemonSet, opts metav1.CreateOptions) (*appsv1.DaemonSet, error)
	GetDaemonSet(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.DaemonSet, error)
	DeleteDaemonSet(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error
	PatchDaemonSet(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*appsv1.DaemonSet, error)
	GetService(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Service, error)
	CreateService(ctx context.Context, namespace string, service *corev1.Service, opts metav1.CreateOptions) (*corev1.Service, error)
	DeleteService(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error
	DeleteDeployment(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error
	CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment, opts metav1.CreateOptions) (*appsv1.Deployment, error)
	GetDeployment(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.Deployment, error)
	PatchDeployment(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*appsv1.Deployment, error)
	CheckDeploymentStatus(ctx context.Context, namespace, deployment string) error
	DeleteNamespace(ctx context.Context, namespace string, opts metav1.DeleteOptions) error
	CreateNamespace(ctx context.Context, namespace string, opts metav1.CreateOptions) (*corev1.Namespace, error)
	GetNamespace(ctx context.Context, namespace string, options metav1.GetOptions) (*corev1.Namespace, error)
	ListPods(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.PodList, error)
	DeletePod(ctx context.Context, namespace, name string, options metav1.DeleteOptions) error
	ExecInPod(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, error)
	CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret, opts metav1.CreateOptions) (*corev1.Secret, error)
	UpdateSecret(ctx context.Context, namespace string, secret *corev1.Secret, opts metav1.UpdateOptions) (*corev1.Secret, error)
	DeleteSecret(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error
	GetSecret(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Secret, error)
	PatchSecret(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*corev1.Secret, error)
	CreateResourceQuota(ctx context.Context, namespace string, r *corev1.ResourceQuota, opts metav1.CreateOptions) (*corev1.ResourceQuota, error)
	DeleteResourceQuota(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error
	AutodetectFlavor(ctx context.Context) k8s.Flavor
	ContextName() (name string)
	CiliumStatus(ctx context.Context, namespace, pod string) (*models.StatusResponse, error)
	ListCiliumEndpoints(ctx context.Context, namespace string, opts metav1.ListOptions) (*ciliumv2.CiliumEndpointList, error)
	GetRunningCiliumVersion(ctx context.Context, namespace string) (string, error)
	GetPlatform(ctx context.Context) (*k8s.Platform, error)
	GetServerVersion() (*semver.Version, error)
	CreateIngressClass(ctx context.Context, r *networkingv1.IngressClass, opts metav1.CreateOptions) (*networkingv1.IngressClass, error)
	DeleteIngressClass(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

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
