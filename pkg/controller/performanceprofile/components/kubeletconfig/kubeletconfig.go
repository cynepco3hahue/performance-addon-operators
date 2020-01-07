package kubeletconfig

import (
	"time"

	performancev1alpha1 "github.com/openshift-kni/performance-addon-operators/pkg/apis/performance/v1alpha1"
	"github.com/openshift-kni/performance-addon-operators/pkg/controller/performanceprofile/components"
	machineconfigv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
)

const (
	cpuManagerPolicyStatic          = "static"
	defaultKubeReservedCPU          = "1000m"
	defaultKubeReservedMemory       = "500Mi"
	defaultSystemReservedCPU        = "1000m"
	defaultSystemReservedMemory     = "500Mi"
	topologyManagerPolicyBestEffort = "best-effort"
)

// NewPerformance returns new KubeletConfig object for performance sensetive workflows
func NewPerformance(profile *performancev1alpha1.PerformanceProfile) *machineconfigv1.KubeletConfig {
	name := components.GetComponentName(profile.Name, components.RoleWorkerPerformance)
	kubeletConfig := &kubeletconfigv1beta1.KubeletConfiguration{
		CPUManagerPolicy:          cpuManagerPolicyStatic,
		CPUManagerReconcilePeriod: metav1.Duration{Duration: 5 * time.Second},
		TopologyManagerPolicy:     topologyManagerPolicyBestEffort,
		KubeReserved: map[string]string{
			"cpu":    defaultKubeReservedCPU,
			"memory": defaultKubeReservedMemory,
		},
		SystemReserved: map[string]string{
			"cpu":    defaultSystemReservedCPU,
			"memory": defaultSystemReservedMemory,
		},
	}

	if profile.Spec.CPU != nil && profile.Spec.CPU.Reserved != nil {
		kubeletConfig.ReservedSystemCPUs = string(*profile.Spec.CPU.Reserved)
	}

	return &machineconfigv1.KubeletConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: machineconfigv1.GroupVersion.String(),
			Kind:       "KubeletConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: machineconfigv1.KubeletConfigSpec{
			MachineConfigPoolSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					components.LabelMachineConfigPoolRole: name,
				},
			},
			KubeletConfig: &runtime.RawExtension{
				Object: kubeletConfig,
			},
		},
	}
}