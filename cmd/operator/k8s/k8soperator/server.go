package k8soperator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// UpdateNodeStatus returns the up to date pod status for each node
func UpdateNodeStatus() map[string]string {
	var status map[string]string
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	options := metav1.ListOptions{FieldSelector: "metadata.name=kubernetes"}
	pods, _ := clientset.CoreV1().Pods("inspr").List(options)

	for _, pod := range pods.Items {
		status[pod.Name] = pod.Status.Message
	}
	return status
}
