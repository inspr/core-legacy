package nodes

import (
	"github.com/inspr/inspr/pkg/meta"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type applyable interface {
	create(no *NodeOperator) error
	del(no *NodeOperator) error
	update(no *NodeOperator) error
}

type kubeService corev1.Service

func (k *kubeService) create(no *NodeOperator) error {
	_, err := no.Services().Create((*corev1.Service)(k))
	return err
}
func (k *kubeService) update(no *NodeOperator) error {
	_, err := no.Services().Update((*corev1.Service)(k))
	return err
}
func (k *kubeService) del(no *NodeOperator) error {
	err := no.Services().Delete(k.Name, nil)
	return err
}

type kubeSecret corev1.Secret

func (k *kubeSecret) create(no *NodeOperator) error {
	_, err := no.Secrets().Create((*corev1.Secret)(k))
	return err
}
func (k *kubeSecret) update(no *NodeOperator) error {
	_, err := no.Secrets().Update((*corev1.Secret)(k))
	return err
}
func (k *kubeSecret) del(no *NodeOperator) error {
	err := no.Secrets().Delete(k.Name, nil)
	return err
}

type kubeDeploy appsv1.Deployment

func (k *kubeDeploy) create(no *NodeOperator) error {
	_, err := no.Deployments().Create((*appsv1.Deployment)(k))
	return err
}
func (k *kubeDeploy) update(no *NodeOperator) error {
	_, err := no.Deployments().Update((*appsv1.Deployment)(k))
	return err
}
func (k *kubeDeploy) del(no *NodeOperator) error {
	err := no.Deployments().Delete(k.Name, nil)
	return err
}

func (no *NodeOperator) dappApplications(app *meta.App) []applyable {
	return []applyable{
		no.dAppToDeployment(app),
		dappToService(app),
	}
}
