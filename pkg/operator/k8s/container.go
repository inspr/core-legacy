package k8s

import corev1 "k8s.io/api/core/v1"

// ContainerOption is a type that changes a container's configuration on instantiation
type ContainerOption func(*corev1.Container)

// ContainerWithCommand changes the command for the container
func ContainerWithCommand(command ...string) ContainerOption {
	return func(c *corev1.Container) {
		c.Command = command
	}
}

// ContainerWithArgs adds arguments to the entrypoint of the container
func ContainerWithArgs(args ...string) ContainerOption {
	return func(c *corev1.Container) {
		c.Args = args
	}
}

// ContainerCWD changes the CWD of a container
func ContainerCWD(cwd string) ContainerOption {
	return func(c *corev1.Container) {
		c.WorkingDir = cwd
	}
}

// ContainerWithPorts adds ports to a k8s container
func ContainerWithPorts(ports ...corev1.ContainerPort) ContainerOption {
	return func(c *corev1.Container) {
		c.Ports = append(c.Ports, ports...)
	}
}

// ContainerWithEnv adds environment variables to a container
func ContainerWithEnv(env ...corev1.EnvVar) ContainerOption {
	return func(c *corev1.Container) {
		c.Env = append(c.Env, env...)
	}
}

// ContainerWithEnvFrom adds environment variables from a source to a container
func ContainerWithEnvFrom(env ...corev1.EnvFromSource) ContainerOption {
	return func(c *corev1.Container) {
		c.EnvFrom = append(c.EnvFrom, env...)
	}
}

// ContainerWithLivenessProbe adds a liveness probe to a k8s container
func ContainerWithLivenessProbe(probe *corev1.Probe) ContainerOption {
	return func(c *corev1.Container) {
		c.LivenessProbe = probe
	}
}

// ContainerWithReadinessProbe adds a readiness probe to a k8s container
func ContainerWithReadinessProbe(probe *corev1.Probe) ContainerOption {
	return func(c *corev1.Container) {
		c.ReadinessProbe = probe
	}
}

// ContainerWithPullPolicy adds a pull policy to a container
func ContainerWithPullPolicy(policy corev1.PullPolicy) ContainerOption {
	return func(c *corev1.Container) {
		c.ImagePullPolicy = policy
	}
}

// WithVolumeMounts adds volume mounts to a container
func WithVolumeMounts(mounts ...corev1.VolumeMount) ContainerOption {
	return func(c *corev1.Container) {
		c.VolumeMounts = append(c.VolumeMounts, mounts...)
	}
}

// NewContainer creates a new kubernetes containers with the given options applied to it
func NewContainer(
	name, image string,
	options ...ContainerOption,
) corev1.Container {
	cont := corev1.Container{
		Name:  name,
		Image: image,
	}
	for _, option := range options {
		if option != nil {
			option(&cont)
		}
	}
	return cont
}
