package builder

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NewService returns a kubernetes service ready to use
func NewService(name string, namespace string, selector map[string]string,
	port int, TargetPort int) apiv1.Service {
	return apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: apiv1.ServiceSpec{
			Selector: selector,
			Ports: []apiv1.ServicePort{
				{
					Port:       int32(port),
					TargetPort: intstr.FromInt(TargetPort),
				},
			},
		},
	}
}
