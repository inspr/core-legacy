package builder

import (
	"reflect"
	"testing"
)

func TestNewPod(t *testing.T) {
	deployment := NewPod()
	if deployment.GetDeployment() == nil {
		t.Errorf("NewPod() returns an empty object.")
	}
}

func Test_builderPod_WithObjectMetadata(t *testing.T) {
	deployment := NewPod()
	labels := map[string]string{
		"app": "inspr-apps-main",
	}

	deployment.WithObjectMetadata("inspr-pod", "inspr-ns", labels)
	deploymentName := deployment.GetDeployment().ObjectMeta.Name
	if deploymentName != "inspr-pod" {
		t.Errorf("NewPod() returns %v expected inspr-pod.", deploymentName)
	}

	deploymentNamespace := deployment.GetDeployment().ObjectMeta.Namespace
	if deploymentNamespace != "inspr-ns" {
		t.Errorf("NewPod() returns %v expected inspr-ns.", deploymentNamespace)
	}

	expectedLabels := map[string]string{
		"app": "inspr-apps-main",
	}
	deploymentLabels := deployment.GetDeployment().ObjectMeta.Labels
	if _, ok := deploymentLabels["app"]; !ok ||
		deploymentLabels["app"] != expectedLabels["app"] {
		t.Errorf("NewPod() returns %v expected inspr-apps-main.",
			deploymentLabels["app"])
	}
}

func Test_builderPod_WithPodTemplateSpec(t *testing.T) {
	deployment := NewPod()
	deployment.WithPodTemplateSpec("inspr-cont", "gcr://inspr.inspr:test")
	deploymentContainer := deployment.GetDeployment().Spec.Template.Spec.Containers[0]

	containerName := deploymentContainer.Name
	if containerName != "inspr-cont" {
		t.Errorf("NewPod() returns %v expected inspr-cont.", containerName)
	}

	containerImage := deploymentContainer.Image
	if containerImage != "gcr://inspr.inspr:test" {
		t.Errorf("NewPod() returns %v expected gcr://inspr.inspr:test.", containerImage)
	}
}

func Test_builderPod_WithPodTemplateObjectMetadata(t *testing.T) {
	deployment := NewPod()
	templateLabels := map[string]string{
		"app": "inspr-apps-main",
	}

	deployment.WithPodTemplateObjectMetadata("inspr-templ", "inspr-ns", templateLabels)
	templateName := deployment.GetDeployment().Spec.Template.ObjectMeta.Name
	if templateName != "inspr-templ" {
		t.Errorf("NewPod() returns %v expected inspr-templ.", templateName)
	}

	templateNamespace := deployment.GetDeployment().Spec.Template.ObjectMeta.Namespace
	if templateNamespace != "inspr-ns" {
		t.Errorf("NewPod() returns %v expected inspr-ns.", templateNamespace)
	}

	expectedLabels := map[string]string{
		"app": "inspr-apps-main",
	}
	deploymentTemplateLabels := deployment.GetDeployment().Spec.Template.ObjectMeta.Labels
	if _, ok := deploymentTemplateLabels["app"]; !ok ||
		deploymentTemplateLabels["app"] != expectedLabels["app"] {
		t.Errorf("NewPod() returns %v expected inspr-apps-main.",
			deploymentTemplateLabels["app"])
	}
}

func Test_builderPod_WithPodSelectorMatchLabels(t *testing.T) {
	deployment := NewPod()
	labels := map[string]string{
		"app": "inspr-apps-main",
	}
	deployment.WithPodSelectorMatchLabels(labels)

	expectedLabels := map[string]string{
		"app": "inspr-apps-main",
	}
	deploymentPodSelectorLabels := deployment.GetDeployment().Spec.Selector.MatchLabels
	if _, ok := deploymentPodSelectorLabels["app"]; !ok ||
		deploymentPodSelectorLabels["app"] != expectedLabels["app"] {
		t.Errorf("NewPod() returns %v expected inspr-apps-main.",
			deploymentPodSelectorLabels["app"])
	}
}

func Test_builderPod_GetDeployment(t *testing.T) {
	deployment := NewPod()
	deploymentObj := deployment.GetDeployment()
	deploymentType := reflect.TypeOf(deploymentObj).String()
	expectedType := "*v1.Deployment"
	if deploymentType != expectedType {
		t.Errorf("NewPod() returns %v expected %v.", deploymentType,
			expectedType)
	}
}
