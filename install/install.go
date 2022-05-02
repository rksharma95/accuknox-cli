package install

import (
	"context"
	"fmt"
	"strings"

	"github.com/kubearmor/kubearmor-client/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

type Options struct {
	Namespace string
}

var selectorLabels = map[string]string{
	"container": "knoxautopolicy",
}

var serviceLabels = map[string]string{
	"service": "knoxautopolicy",
}

var deploymentLabels = map[string]string{
	"deployment": "knoxautopolicy",
}

func GetService(namespace string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "knoxautopolicy",
			Labels: serviceLabels,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Port:       9089,
					TargetPort: intstr.FromInt(9089),
					Protocol:   "TCP",
				},
			},
		},
	}
}

func GetDeployment(namespace string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "knoxautopolicy",
			Labels:    deploymentLabels,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: selectorLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "knoxautopolicy",
					Containers: []corev1.Container{
						{
							Name:            "knoxautopolicy",
							Image:           "accuknox/knoxautopolicy:dev",
							ImagePullPolicy: "Always",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9089,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config-volume", //BPF (read-only)
									MountPath: "/config",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "knoxautopolicy-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func GetServiceAccount(namespace string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "knoxautopolicy",
			Namespace: namespace,
		},
	}
}

func GetClusterRoleBinding(namespace string) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "knoxautopolicy",
			Namespace: namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "knoxautopolicy",
				Namespace: namespace,
			},
		},
	}
}

func DiscoveryEngineInstaller(c *k8s.Client, o Options) error {

	o.Namespace = "explorer"

	// discovery-engine Service
	fmt.Print("Discovery-engine Service...\n")
	if _, err := c.K8sClientset.CoreV1().Services(o.Namespace).Create(context.Background(), GetService(o.Namespace), metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
		fmt.Print("Discovery-engine Service already exists...\n")
	}

	//TODO discovery-engine dev-config
	// fmt.Print("KubeArmor Relay Service ...\n")
	// if _, err := c.K8sClientset.CoreV1().ConfigMaps(o.Namespace).Create(&de_configmap); err != nil {
	// 	if !strings.Contains(err.Error(), "already exists") {
	// 		return err
	// 	}
	// 	fmt.Print("KubeArmor Relay Service already exists ...\n")
	// }

	// discovery-engine Deployment
	fmt.Print("KubeArmor Relay Deployment...\n")
	if _, err := c.K8sClientset.AppsV1().Deployments(o.Namespace).Create(context.Background(), GetDeployment(o.Namespace), metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
		fmt.Print("KubeArmor Relay Deployment already exists...\n")
	}

	// discovery-engine Service account
	fmt.Print("Service Account...\n")
	if _, err := c.K8sClientset.CoreV1().ServiceAccounts(o.Namespace).Create(context.Background(), GetServiceAccount(o.Namespace), metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
		fmt.Print("Service Account already exists...\n")
	}

	// discovery-engine cluster role bindings
	fmt.Print("Cluster Role Bindings...\n")
	if _, err := c.K8sClientset.RbacV1().ClusterRoleBindings().Create(context.Background(), GetClusterRoleBinding(o.Namespace), metav1.CreateOptions{}); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
		fmt.Print("Cluster Role Bindings already exists...\n")
	}

	return nil
}

func DiscoveryEngineUninstaller(c *k8s.Client, o Options) error {

	o.Namespace = "explorer"

	// discovery-engine Service
	fmt.Print("🔥 Deleting Discovery-engine Service...\n")
	if err := c.K8sClientset.CoreV1().Services(o.Namespace).Delete(context.Background(), "knoxautopolicy", metav1.DeleteOptions{}); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
		fmt.Print("Discovery-engine Service not found...\n")
	}

	//TODO discovery-engine dev-config
	// fmt.Print("🔥 Deleting KubeArmor Relay Service...\n")
	// if err := c.K8sClientset.CoreV1().ConfigMaps(o.Namespace).Delete(&de_configmap); err != nil {
	// 	if !strings.Contains(err.Error(), "not found") {
	// 		return err
	// 	}
	// 	fmt.Print("KubeArmor Relay Service not found...\n")
	// }

	// discovery-engine Deployment
	fmt.Print("🔥 Deleting KubeArmor Relay Deployment...\n")
	if err := c.K8sClientset.AppsV1().Deployments(o.Namespace).Delete(context.Background(), "knoxautopolicy", metav1.DeleteOptions{}); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
		fmt.Print("KubeArmor Relay Deployment not found...\n")
	}

	// discovery-engine Service account
	fmt.Print("🔥 Deleting Service Account...\n")
	if err := c.K8sClientset.CoreV1().ServiceAccounts(o.Namespace).Delete(context.Background(), "knoxautopolicy", metav1.DeleteOptions{}); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
		fmt.Print("Service Account not found...\n")
	}

	// discovery-engine cluster role bindings
	fmt.Print("🔥 Deleting Cluster Role Bindings...\n")
	if err := c.K8sClientset.RbacV1().ClusterRoleBindings().Delete(context.Background(), "knoxautopolicy", metav1.DeleteOptions{}); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
		fmt.Print("Cluster Role Bindings not found...\n")
	}

	return nil
}
