package builder

import (
	"reflect"
	"strconv"
	"testing"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_builderIngress_NewIngress(t *testing.T) {
	ingress := NewIngress()
	if ingress.GetIngress() == nil {
		t.Errorf("NewIngress() returns an empty object.")
	}
}

func Test_builderIngress_GetIngress(t *testing.T) {
	ingress := NewIngress()
	ingressObj := ingress.GetIngress()
	ingressType := reflect.TypeOf(ingressObj).String()
	expectedType := "*v1beta1.Ingress"
	if ingressType != expectedType {
		t.Errorf("NewIngress() returns %v expected %v.", ingressType,
			expectedType)
	}
}

func Test_builderIngress_WithIngressBackend(t *testing.T) {
	ingress := NewIngress()
	ingress.WithIngressBackend("test-svc", "91")
	ingressObj := ingress.GetIngress()

	backendServiceName := ingressObj.Spec.Backend.ServiceName
	if backendServiceName != "test-svc" {
		t.Errorf("NewIngress() returns %v expected test-svc.", backendServiceName)
	}

	backendServicePort := ingressObj.Spec.Backend.ServicePort
	expectedPort := intstr.FromInt(91)
	if backendServicePort != expectedPort {
		t.Errorf("NewIngress() returns %v expected %v.", backendServicePort, expectedPort)
	}

	ingress.WithIngressBackend("test2-svc", "")

	backendServicePort = ingressObj.Spec.Backend.ServicePort
	if backendServicePort.Type != 0 {
		t.Errorf("NewIngress() returns %#v expected nil.", backendServicePort)
	}

	backendServiceName = ingressObj.Spec.Backend.ServiceName
	if backendServiceName != "test2-svc" {
		t.Errorf("NewIngress() returns %v expected test-svc.", backendServiceName)
	}
}

func Test_builderIngress_WithIngressRules(t *testing.T) {
	ingress := NewIngress()
	ingress.WithIngressRules("www.inspr.com", "inspr-svc", "87")
	ingressObj := ingress.GetIngress()

	ingressHost := ingressObj.Spec.Rules[0].Host

	if ingressHost != "www.inspr.com" {
		t.Errorf("NewIngress() returns %v expected www.inspr.com.",
			ingressHost)
	}

	ingressServiceName := ingressObj.Spec.Backend.ServiceName
	if ingressServiceName != "inspr-svc" {
		t.Errorf("NewIngress() returns %v expected test-svc.", ingressServiceName)
	}

	ingressServicePort := ingressObj.Spec.Backend.ServicePort
	if ingressServicePort != intstr.FromInt(87) {
		t.Errorf("NewIngress() returns %v expected 91.", ingressServicePort)
	}
}

func Test_builderIngress_WithObjectMeta(t *testing.T) {
	annotations := map[string]string{
		"kubernetes.io/ingress.class": "nginx",
	}

	ingress := NewIngress()
	ingress.WithObjectMeta("inspr-ingress", "inspr-namespace",
		annotations)
	ingressObj := ingress.GetIngress()

	metadataName := ingressObj.ObjectMeta.Name
	if metadataName != "inspr-ingress" {
		t.Errorf("NewIngress() returns %v expected inspr-ingress.",
			metadataName)
	}

	metadataNamespace := ingressObj.ObjectMeta.Namespace
	if metadataNamespace != "inspr-namespace" {
		t.Errorf("NewIngress() returns %v expected inspr-namespace.",
			metadataNamespace)
	}

	ingressAnnotations := ingressObj.ObjectMeta.Annotations
	if ingressAnnotations["kubernetes.io/ingress.class"] != "nginx" {
		t.Errorf("NewIngress() returns %v expected nginx.",
			ingressAnnotations["kubernetes.io/ingress.class"])
	}

	annotations = map[string]string{}
	ingress.WithObjectMeta("inspr-ingress-2", "inspr-namespace-2",
		annotations)

	metadataName = ingressObj.ObjectMeta.Name
	if metadataName != "inspr-ingress-2" {
		t.Errorf("NewIngress() returns %v expected inspr-ingress-2.",
			metadataName)
	}

	metadataNamespace = ingressObj.ObjectMeta.Namespace
	if metadataNamespace != "inspr-namespace-2" {
		t.Errorf("NewIngress() returns %v expected inspr-namespace-2.",
			metadataNamespace)
	}

	ingressAnnotations = ingressObj.ObjectMeta.Annotations
	if len(ingressAnnotations) != 0 {
		t.Errorf("NewIngress() returns %#v expected nil.",
			ingressAnnotations)
	}

}

func Test_builderIngress_WithTLS(t *testing.T) {
	ingress := NewIngress()
	ingress.WithTLS("inspr.", "company.domain", "inspr")
	ingressObj := ingress.GetIngress()

	host := ingressObj.Spec.TLS[0].Hosts[0]
	if host != "inspr.company.domain" {
		t.Errorf("NewIngress() returns %v expected inspr.company.domain", host)
	}

	secretName := ingressObj.Spec.TLS[0].SecretName
	if secretName != "inspr-secret" {
		t.Errorf("NewIngress() returns %v expected inspr-namespaces.",
			secretName)
	}
}

func TestNewFromIngress(t *testing.T) {
	var ing *v1beta1.Ingress
	ingress := NewFromIngress(ing)
	if ingress == nil {
		t.Errorf("NewFromIngress() returns nil expected %#v.", ingress)
	}
}

func Test_remove(t *testing.T) {
	var paths []v1beta1.HTTPIngressPath
	var path v1beta1.HTTPIngressPath

	for index := 0; index < 5; index++ {
		path.Path = "test" + strconv.Itoa(index)
		paths = append(paths, path)
	}

	paths = remove(paths, 3)

	for row := range paths {
		if paths[row].Path == "test3" {
			t.Errorf("remove() did not delete %#v.", "test3")
		}
	}
}

func Test_builderIngress_WithRule(t *testing.T) {
	ingress := NewIngress()
	ingress.WithIngressRules("inspr.inspr.com", "main-svc", "37")
	ingress.WithRule("inspr-svc", 37, "/test")

	if len(ingress.GetIngress().Spec.Rules) != 1 {
		t.Errorf("NewIngress().WithRule() did not create the rule.")
	}

	if len(ingress.GetIngress().Spec.Rules[0].HTTP.Paths) != 2 {
		t.Errorf("NewIngress().WithRule() did not create the rule.")
	}

	path := ingress.GetIngress().Spec.Rules[0].HTTP.Paths[1].Path
	if path != "/test" {
		t.Errorf("NewIngress().WithRule() returns %#v expected %#v.",
			path, "/test")
	}

	serviceName := ingress.GetIngress().Spec.Rules[0].HTTP.Paths[1].Backend.ServiceName
	if serviceName != "inspr-svc" {
		t.Errorf("NewIngress().WithRule() returns %#v expected %#v.",
			serviceName, "inspr-svc")
	}

	ServicePort := ingress.GetIngress().Spec.Rules[0].HTTP.Paths[1].Backend.ServicePort
	if ServicePort != intstr.FromInt(37) {
		t.Errorf("NewIngress().WithRule() returns %#v expected %#v.",
			ServicePort, intstr.FromInt(37))
	}

	ingress.WithRule("inspr-svc", 42, "/test")

	path = ingress.GetIngress().Spec.Rules[0].HTTP.Paths[1].Path
	if path != "/test" {
		t.Errorf("NewIngress().WithRule() returns %#v expected %#v.",
			path, "/test")
	}

	serviceName = ingress.GetIngress().Spec.Rules[0].HTTP.Paths[1].Backend.ServiceName
	if serviceName != "inspr-svc" {
		t.Errorf("NewIngress().WithRule() returns %#v expected %#v.",
			serviceName, "inspr-svc")
	}

	ServicePort = ingress.GetIngress().Spec.Rules[0].HTTP.Paths[1].Backend.ServicePort
	if ServicePort != intstr.FromInt(42) {
		t.Errorf("NewIngress().WithRule() returns %#v expected %#v.",
			ServicePort, intstr.FromInt(42))
	}
}
