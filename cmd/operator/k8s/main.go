package main

import (
	k8spackop "gitlab.inspr.dev/inspr/core/cmd/operator/k8s/k8spack"
)

func main() {
	url := ""
	port := "50000"
	k8spackop.NewK8sOperator(url, port)
}
