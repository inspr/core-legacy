package builder

import (
	"strconv"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Ingress is used to build an Ingress.
type Ingress interface {
	WithIngressBackend(serviceName string, servicePort string) Ingress
	WithIngressRules(host string, serviceName string, servicePort string) Ingress
	WithObjectMeta(name string, namespace string, annotations map[string]string) Ingress
	WithTLS(subdomain string, domain string, appsNamespace string) Ingress
	WithRule(serviceName string, servicePort int, path string) Ingress
	GetIngress() *v1beta1.Ingress
}

type builderIngress struct {
	ingress *v1beta1.Ingress
}

// NewIngress creates a new command builder
func NewIngress() Ingress {
	return &builderIngress{
		ingress: &v1beta1.Ingress{},
	}
}

// NewFromIngress creates an ingress builder from a pre existing ingress
func NewFromIngress(ing *v1beta1.Ingress) Ingress {
	return &builderIngress{
		ingress: ing,
	}
}

// GetIngress returns the ingress struct object
func (builder *builderIngress) GetIngress() *v1beta1.Ingress {
	return builder.ingress
}

// WithIngressBackend Creates an Ingress with Backend based
// in the parameters
func (builder *builderIngress) WithIngressBackend(serviceName string,
	servicePort string) Ingress {
	if servicePort == "" {
		builder.ingress.Spec.Backend = &v1beta1.IngressBackend{
			ServiceName: serviceName,
		}
	} else {
		port, _ := strconv.Atoi(servicePort)
		builder.ingress.Spec.Backend = &v1beta1.IngressBackend{
			ServiceName: serviceName,
			ServicePort: intstr.FromInt(port),
		}
	}
	return builder
}

// WithIngressRules Creates an Ingress with Rules based
// in the parameters
func (builder *builderIngress) WithIngressRules(host string, serviceName string,
	servicePort string) Ingress {
	builder.WithIngressBackend(serviceName, servicePort)
	builder.ingress.Spec.Rules = []v1beta1.IngressRule{
		{
			Host: host,
			IngressRuleValue: v1beta1.IngressRuleValue{
				HTTP: &v1beta1.HTTPIngressRuleValue{
					Paths: []v1beta1.HTTPIngressPath{
						{
							Path:    "/",
							Backend: *builder.ingress.Spec.Backend,
						},
					},
				},
			},
		},
	}
	return builder
}

// WithObjectMeta Creates an Ingress with an Object metadata based
// in the parameters
func (builder *builderIngress) WithObjectMeta(name string, namespace string,
	annotations map[string]string) Ingress {
	if len(annotations) == 0 {
		builder.ingress.ObjectMeta = metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		}
	}
	builder.ingress.ObjectMeta = metav1.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Annotations: annotations,
	}
	return builder
}

// WithTLS Creates an Ingress with a TLS metadata based
// in the parameters
func (builder *builderIngress) WithTLS(subdomain string, domain string,
	appsNamespace string) Ingress {
	builder.ingress.Spec.TLS = []v1beta1.IngressTLS{
		{
			Hosts:      []string{subdomain + domain},
			SecretName: appsNamespace + "-secret",
		},
	}
	return builder
}
func remove(s []v1beta1.HTTPIngressPath, i int) []v1beta1.HTTPIngressPath {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (builder *builderIngress) WithRule(serviceName string, servicePort int, path string) Ingress {
	ingress := builder.GetIngress()

	paths := ingress.Spec.Rules[0].HTTP.Paths

	for i := range paths {
		if paths[i].Path == path {
			ingress.Spec.Rules[0].HTTP.Paths = remove(paths, i)
			break
		}
	}

	ingress.Spec.Rules[0].HTTP.Paths = append(
		ingress.Spec.Rules[0].HTTP.Paths,
		v1beta1.HTTPIngressPath{
			Path: path,
			Backend: v1beta1.IngressBackend{
				ServiceName: serviceName,
				ServicePort: intstr.FromInt(servicePort),
			},
		})
	builder.ingress = ingress
	return builder
}
