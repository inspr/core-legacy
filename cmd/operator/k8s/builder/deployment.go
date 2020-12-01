package builder

import (
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod is used to build a pod with a deployment.
type Pod interface {
	WithObjectMetadata(name string, namespace string, labels map[string]string) Pod
	WithPodTemplateSpec(containerName string, containerImage string) Pod
	WithPodTemplateObjectMetadata(
		templateName string,
		templateNamespace string,
		templateLabels map[string]string,
	) Pod
	WithPodSelectorMatchLabels(matchLabels map[string]string) Pod
	GetDeployment() *v1.Deployment
}

type builderPod struct {
	pod v1.Deployment
}

// NewPod creates a new command builder
func NewPod() Pod {
	return &builderPod{
		pod: v1.Deployment{},
	}
}

//WithObejctMetadata return a Object Metadata to a pod
func (builder *builderPod) WithObjectMetadata(name string, namespace string,
	labels map[string]string) Pod {
	builder.pod.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
	return builder
}

// WithPodSpec specify the pod spec to a pod template.
func (builder *builderPod) WithPodTemplateSpec(containerName string, containerImage string) Pod {
	builder.pod.Spec.Template.Spec = apiv1.PodSpec{
		Containers: []apiv1.Container{
			{
				Name:  containerName,
				Image: containerImage,
				// parse from master env var to kube env vars
				ImagePullPolicy: apiv1.PullAlways,
			},
		},
	}
	return builder
}

// WithObjectMetadata specify the metadata to a pod template.
func (builder *builderPod) WithPodTemplateObjectMetadata(templateName string, templateNamespace string,
	templateLabels map[string]string) Pod {
	builder.pod.Spec.Template.ObjectMeta = metav1.ObjectMeta{
		Name:      templateName,
		Namespace: templateNamespace,
		Labels:    templateLabels,
	}
	return builder
}

// WithPodSelectorMatchLabels create a Match label for a given pod.
func (builder *builderPod) WithPodSelectorMatchLabels(matchLabels map[string]string) Pod {
	builder.pod.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: matchLabels,
	}
	return builder
}

// GetDeployment returns the created object
func (builder *builderPod) GetDeployment() *v1.Deployment {
	return &builder.pod
}
