package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"math/big"
	"strings"

	"inspr.dev/inspr/pkg/controller/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func generatePassword() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 20
	var b strings.Builder
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b.WriteRune(chars[index.Int64()])
	}
	str := b.String()
	return str
}

var clientSet kubernetes.Interface

func initInsprd() (string, error) {

	cont := client.NewControllerClient(client.ControllerConfig{
		URL: os.Getenv("INSPRD_URL"),
	})

	token, err := cont.Authorization().
		Init(context.Background(), os.Getenv("INSPRD_INIT_KEY"))
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
	secret, err := clientSet.CoreV1().
		Secrets(namespace).
		Get(ctx, secretName, v1.GetOptions{})
	if err != nil {
		panic(err)
	}

	if _, exists := secret.Data["REFRESH_KEY"]; !exists {
		bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
		if _, err := rand.Read(bytes); err != nil {
			panic(err.Error())
		}

		key := hex.EncodeToString(bytes)
		privateKeyBytes := []byte(key)
		secret.Data["REFRESH_KEY"] = privateKeyBytes
		if os.Getenv("INIT_INSPRD") == "true" {
			token, err := initInsprd()
			if err != nil {
				panic(err)
			}
			secret.Data["ADMIN_TOKEN"] = []byte(token)
		}
		if os.Getenv("ADMIN_PASSWORD_GENERATE") == "true" {
			secret.Data["ADMIN_PASSWORD"] = []byte(generatePassword())
		}

		_, err = clientSet.CoreV1().
			Secrets(namespace).
			Update(ctx, secret, v1.UpdateOptions{})
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
