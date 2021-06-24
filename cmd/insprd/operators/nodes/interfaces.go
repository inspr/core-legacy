package nodes

import (
	"context"
	"log"

	"inspr.dev/inspr/pkg/meta"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type applyable interface {
	create(no *NodeOperator) error
	del(no *NodeOperator) error
	update(no *NodeOperator) error
}

type kubeService corev1.Service

func (k *kubeService) create(no *NodeOperator) error {
	_, err := no.Services().Create(context.Background(), (*corev1.Service)(k), v1.CreateOptions{})
	return err
}
func (k *kubeService) update(no *NodeOperator) error {
	_, err := no.Services().Update(context.Background(), (*corev1.Service)(k), v1.UpdateOptions{})
	return err
}
func (k *kubeService) del(no *NodeOperator) error {
	err := no.Services().Delete(context.Background(), k.Name, v1.DeleteOptions{})
	return err
}

type kubeSecret corev1.Secret

func (k *kubeSecret) create(no *NodeOperator) error {
	log.Println("creating secret")
	_, err := no.Secrets().Create(context.Background(), (*corev1.Secret)(k), v1.CreateOptions{})
	return err
}
func (k *kubeSecret) update(no *NodeOperator) error {
	_, err := no.Secrets().Update(context.Background(), (*corev1.Secret)(k), v1.UpdateOptions{})
	return err
}
func (k *kubeSecret) del(no *NodeOperator) error {
	err := no.Secrets().Delete(context.Background(), k.Name, v1.DeleteOptions{})
	return err
}

type kubeDeploy appsv1.Deployment

func (k *kubeDeploy) create(no *NodeOperator) error {
	_, err := no.Deployments().Create(context.Background(), (*appsv1.Deployment)(k), v1.CreateOptions{})
	return err
}
func (k *kubeDeploy) update(no *NodeOperator) error {
	_, err := no.Deployments().Update(context.Background(), (*appsv1.Deployment)(k), v1.UpdateOptions{})
	return err
}
func (k *kubeDeploy) del(no *NodeOperator) error {
	err := no.Deployments().Delete(context.Background(), k.Name, v1.DeleteOptions{})
	return err
}

func (no *NodeOperator) dappApplications(app *meta.App) []applyable {
	return []applyable{
		no.toSecret(app),
		no.dAppToDeployment(app),
		no.dappToService(app),
	}
}
