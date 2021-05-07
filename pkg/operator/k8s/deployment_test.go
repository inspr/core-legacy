package k8s

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWithSelector(t *testing.T) {
	type args struct {
		sel *metav1.LabelSelector
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "adding selector",
			args: args{
				sel: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "this is an app",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "this is an app",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.Deployment{}
			option := WithSelector(tt.args.sel)
			option(dep)
			if !reflect.DeepEqual(dep, tt.want) {
				t.Errorf("WithSelector() = %v, want %v", dep, tt.want)
			}
		})
	}
}

func TestWithContainer(t *testing.T) {
	type args struct {
		cont []corev1.Container
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "multiple containers",
			args: args{
				cont: []corev1.Container{
					{
						Name: "container1",
					},
					{
						Name: "container2",
					},
					{
						Name: "container3",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name: "container1",
								},
								{
									Name: "container2",
								},
								{
									Name: "container3",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "single container",
			args: args{
				cont: []corev1.Container{
					{
						Name: "container1",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name: "container1",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.Deployment{}
			option := WithContainer(tt.args.cont...)
			option(dep)
			if !reflect.DeepEqual(dep, tt.want) {
				t.Errorf("WithContainer() = %v, want %v", dep, tt.want)
			}
		})
	}
}

func TestWithInitContainers(t *testing.T) {
	type args struct {
		cont []corev1.Container
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "multiple containers",
			args: args{
				cont: []corev1.Container{
					{
						Name: "container1",
					},
					{
						Name: "container2",
					},
					{
						Name: "container3",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							InitContainers: []corev1.Container{
								{
									Name: "container1",
								},
								{
									Name: "container2",
								},
								{
									Name: "container3",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "single container",
			args: args{
				cont: []corev1.Container{
					{
						Name: "container1",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							InitContainers: []corev1.Container{
								{
									Name: "container1",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.Deployment{}
			option := WithInitContainers(tt.args.cont...)
			option(dep)
			if !reflect.DeepEqual(dep, tt.want) {
				t.Errorf("WithContainer() = %v, want %v", dep, tt.want)
			}
		})
	}
}

func TestWithVolumes(t *testing.T) {
	type args struct {
		cont []corev1.Volume
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "multiple volumes",
			args: args{
				cont: []corev1.Volume{
					{
						Name: "volume1",
					},
					{
						Name: "volume2",
					},
					{
						Name: "volume3",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								{
									Name: "volume1",
								},
								{
									Name: "volume2",
								},
								{
									Name: "volume3",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "single volume",
			args: args{
				cont: []corev1.Volume{
					{
						Name: "volume1",
					},
				},
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								{
									Name: "volume1",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.Deployment{}
			option := WithVolumes(tt.args.cont...)
			option(dep)
			if !reflect.DeepEqual(dep, tt.want) {
				t.Errorf("WithVolume() = %v, want %v", dep, tt.want)
			}
		})
	}
}

func TestWithReplicas(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "correct injection",
			args: args{
				n: 13,
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Replicas: func() *int32 {
						var k int32 = 13
						return &k
					}(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.Deployment{}
			option := WithReplicas(tt.args.n)
			option(dep)
			if !reflect.DeepEqual(dep, tt.want) {
				t.Errorf("WithVolume() = %v, want %v", dep, tt.want)
			}
		})
	}
}

func TestWithLabels(t *testing.T) {
	type args struct {
		labels map[string]string
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "lots of labels",
			args: args{
				labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"label1": "value1",
							"label2": "value2",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"label1": "value1",
								"label2": "value2",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.Deployment{}
			option := WithLabels(tt.args.labels)
			option(dep)
			if !reflect.DeepEqual(dep, tt.want) {
				t.Errorf("WithVolume() = %v, want %v", dep, tt.want)
			}
		})
	}
}

func TestWithAnnotations(t *testing.T) {
	type args struct {
		labels map[string]string
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "lots of labels",
			args: args{
				labels: map[string]string{
					"label1": "value1",
					"label2": "value2",
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"label1": "value1",
						"label2": "value2",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								"label1": "value1",
								"label2": "value2",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.Deployment{}
			option := WithAnnotations(tt.args.labels)
			option(dep)
			if !reflect.DeepEqual(dep, tt.want) {
				t.Errorf("WithVolume() = %v, want %v", dep, tt.want)
			}
		})
	}
}

func TestWithRestartPolicy(t *testing.T) {
	type args struct {
		policy corev1.RestartPolicy
	}
	tests := []struct {
		name string
		args args
		want *appsv1.Deployment
	}{
		{
			name: "correct injection",
			args: args{
				policy: corev1.RestartPolicyNever,
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
		{
			name: "correct injection",
			args: args{
				policy: corev1.RestartPolicyAlways,
			},
			want: &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyAlways,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dep := &appsv1.Deployment{}
			option := WithRestartPolicy(tt.args.policy)
			option(dep)
			if !reflect.DeepEqual(dep, tt.want) {
				t.Errorf("WithVolume() = %v, want %v", dep, tt.want)
			}
		})
	}
}
func assertEQ(t *testing.T, f string, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v got = %v, want %v", f, got, want)
	}
}
func TestNewDeployment(t *testing.T) {
	name := "name1"
	ran1 := false
	ran2 := false
	testingOption1 := func(dep *appsv1.Deployment) {
		assertEQ(t, "NewDeployment", dep.Name, name)
		ran1 = true
	}
	testingOption2 := func(dep *appsv1.Deployment) {
		ran2 = true
	}

	dep := NewDeployment(name, testingOption1, testingOption2)
	assertEQ(t, "NewDeployment", dep, &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: name}})

	assertEQ(t, "NewDeployment", ran1, true)
	assertEQ(t, "NewDeployment", ran2, true)
}
