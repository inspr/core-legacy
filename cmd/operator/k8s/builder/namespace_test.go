package builder

import "testing"

func TestNewNamespace(t *testing.T) {
	kubeNamespace := NewNamespace("inspr-test")
	kubeRegistredNamespace := kubeNamespace.GetObjectMeta().GetName()
	if kubeRegistredNamespace != "inspr-test" {
		t.Errorf("NewNamespace() returns %v expected inspr-test.",
			kubeRegistredNamespace)
	}
}
