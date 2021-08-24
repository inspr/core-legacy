package main

import (
	"context"
	"log"
	"os"

	"inspr.dev/inspr/pkg/controller/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientSet kubernetes.Interface

func initInsprd() ( string, error ){

	cont :=  client.NewControllerClient(client.ControllerConfig{
		URL: os.Getenv("INSPRD_URL"),
	})

	token, err := cont.Authorization().Init(context.Background(), os.Getenv("INSPRD_INIT_KEY"))
	return token, err
}
// initKube initializes a k8s operator with in cluster configuration
func initKube() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	namespace := os.Getenv("K8S_NAMESPACE")
	secretName := os.Getenv("SECRET_NAME")

	initKube()
	secret, err := clientSet.CoreV1().Secrets(namespace).Get(ctx,secretName, v1.GetOptions{})
	if err != nil {
		panic(err)
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	token, err := initInsprd()
	if err != nil {
		panic(err)
	}
	secret.Data["ADMIN_TOKEN"] = []byte(token)

	_, err = clientSet.CoreV1().Secrets(namespace).Update(ctx, secret, v1.UpdateOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}
}
