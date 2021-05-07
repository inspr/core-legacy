package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentOption is a type to change a deployment on instantiation
type DeploymentOption func(*appsv1.Deployment)

// WithSelector changes the selector for the deployment
func WithSelector(sel *metav1.LabelSelector) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Selector = sel
	}
}

// WithContainer adds containers to the deployment
func WithContainer(cont ...corev1.Container) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, cont...)
	}
}

// WithInitContainers adds initialization containers to the deployment
func WithInitContainers(cont ...corev1.Container) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Template.Spec.InitContainers = append(d.Spec.Template.Spec.InitContainers, cont...)
	}
}

// WithVolumes adds volumes to a deployment
func WithVolumes(vol ...corev1.Volume) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, vol...)
	}
}

// WithReplicas changes the number of replicas of a deployment
func WithReplicas(n int) DeploymentOption {
	converted := int32(n)
	return func(d *appsv1.Deployment) {
		d.Spec.Replicas = &converted
	}
}

// WithLabels adds labels to the deployment. Also sets the label selector for deployment pods
func WithLabels(labels map[string]string) DeploymentOption {
	return func(d *appsv1.Deployment) {
		if d.Labels == nil {
			d.Labels = make(map[string]string)
		}
		if d.Spec.Template.Labels == nil {
			d.Spec.Template.Labels = make(map[string]string)
		}
		for key, value := range labels {
			d.Labels[key] = value
			d.Spec.Template.Labels[key] = value
		}
		d.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
	}
}

// WithAnnotations adds annotations to the deployment and it's templates
func WithAnnotations(labels map[string]string) DeploymentOption {
	return func(d *appsv1.Deployment) {
		if d.Annotations == nil {
			d.Annotations = make(map[string]string)
		}
		if d.Spec.Template.Annotations == nil {
			d.Spec.Template.Annotations = make(map[string]string)
		}
		for key, value := range labels {
			d.Annotations[key] = value
			d.Spec.Template.Annotations[key] = value
		}
	}
}

// WithRestartPolicy changes the restart policy of the deployment's template
func WithRestartPolicy(policy corev1.RestartPolicy) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Template.Spec.RestartPolicy = policy
	}
}

// NewDeployment creates a deployment with the given options applied to it
func NewDeployment(name string, options ...DeploymentOption) *appsv1.Deployment {
	cont := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	for _, option := range options {
		option(cont)
	}
	return cont
}
