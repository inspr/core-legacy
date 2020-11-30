package builder

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestNewService(t *testing.T) {
	selector := map[string]string{
		"app": "inspr-apps-main",
	}
	kubeService := NewService("inspr-svc", "inspr-apps", selector, 91, 92)

	kubeServiceName := kubeService.GetObjectMeta().GetName()
	if kubeServiceName != "inspr-svc" {
		t.Errorf("NewService() returns %v expected inspr-svc.", kubeServiceName)
	}

	kubeServiceNamespace := kubeService.GetObjectMeta().GetNamespace()
	if kubeServiceNamespace != "inspr-apps" {
		t.Errorf("NewService() returns %v expected inspr-apps.",
			kubeServiceNamespace)
	}

	kubeServiceSelector := kubeService.Spec.Selector
	expectedSelector := map[string]string{
		"app": "inspr-apps-main",
	}
	if _, ok := kubeServiceSelector["app"]; !ok ||
		kubeServiceSelector["app"] != expectedSelector["app"] {
		t.Errorf("NewService() returns %v expected inspr-apps-main.",
			kubeServiceSelector["app"])
	}

	kubeServicePort := kubeService.Spec.Ports[0].Port
	if kubeServicePort != 91 {
		t.Errorf("NewService() returns %v expected 91.", kubeServicePort)
	}

	kubeServiceTargetPort := kubeService.Spec.Ports[0].TargetPort
	if kubeServiceTargetPort != intstr.FromInt(92) {
		t.Errorf("NewService() returns %v expected 92.", kubeServiceTargetPort)
	}
}
