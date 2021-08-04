package nodes

import (
	"context"

	"go.uber.org/zap"
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
	logger.Info("creating service resource on kubernetes", zap.String("service-name", k.ObjectMeta.Name))
	_, err := no.Services().Create(context.Background(), (*corev1.Service)(k), v1.CreateOptions{})
	if err != nil {
		logger.Error("unable to create service resource on kubernetes", zap.String("service-name", k.ObjectMeta.Name))
	}
	return err
}
func (k *kubeService) update(no *NodeOperator) error {
	logger.Info("updating service resource on kubernetes", zap.String("service-name", k.ObjectMeta.Name))
	_, err := no.Services().Update(context.Background(), (*corev1.Service)(k), v1.UpdateOptions{})
	if err != nil {
		logger.Error("unable to update service resource on kubernetes", zap.String("service-name", k.ObjectMeta.Name))
	}
	return err
}
func (k *kubeService) del(no *NodeOperator) error {
	logger.Info("deleting service resource on kubernetes", zap.String("service-name", k.ObjectMeta.Name))
	err := no.Services().Delete(context.Background(), k.Name, v1.DeleteOptions{})
	if err != nil {
		logger.Error("unable to update service resource on kubernetes", zap.String("service-name", k.ObjectMeta.Name))
	}
	return err
}

type kubeSecret corev1.Secret

func (k *kubeSecret) create(no *NodeOperator) error {
	logger.Info("creating secret resource on kubernetes", zap.String("secret-name", k.ObjectMeta.Name))
	_, err := no.Secrets().Create(context.Background(), (*corev1.Secret)(k), v1.CreateOptions{})
	if err != nil {
		logger.Error("unable to create secret resource on kubernetes", zap.String("secret-name", k.ObjectMeta.Name))
	}
	return err
}
func (k *kubeSecret) update(no *NodeOperator) error {
	logger.Info("updating secret resource on kubernetes", zap.String("secret-name", k.ObjectMeta.Name))
	_, err := no.Secrets().Update(context.Background(), (*corev1.Secret)(k), v1.UpdateOptions{})
	if err != nil {
		logger.Error("unable to update secret resource on kubernetes", zap.String("secret-name", k.ObjectMeta.Name))
	}
	return err
}
func (k *kubeSecret) del(no *NodeOperator) error {
	logger.Info("deleting secret resource on kubernetes", zap.String("secret-name", k.ObjectMeta.Name))
	err := no.Secrets().Delete(context.Background(), k.Name, v1.DeleteOptions{})
	if err != nil {
		logger.Error("unable to update secret resource on kubernetes", zap.String("secret-name", k.ObjectMeta.Name))
	}
	return err
}

type kubeDeployment appsv1.Deployment

func (k *kubeDeployment) create(no *NodeOperator) error {
	logger.Info("creating deployment resource on kubernetes", zap.String("deployment-name", k.ObjectMeta.Name))
	_, err := no.Deployments().Create(context.Background(), (*appsv1.Deployment)(k), v1.CreateOptions{})
	if err != nil {
		logger.Error("unable to create deployment resource on kubernetes", zap.String("deployment-name", k.ObjectMeta.Name))
	}
	return err
}
func (k *kubeDeployment) update(no *NodeOperator) error {
	logger.Info("updating deployment resource on kubernetes", zap.String("deployment-name", k.ObjectMeta.Name))
	_, err := no.Deployments().Update(context.Background(), (*appsv1.Deployment)(k), v1.UpdateOptions{})
	if err != nil {
		logger.Error("unable to update deployment resource on kubernetes", zap.String("deployment-name", k.ObjectMeta.Name))
	}
	return err
}
func (k *kubeDeployment) del(no *NodeOperator) error {
	logger.Info("deleting deployment resource on kubernetes", zap.String("deployment-name", k.ObjectMeta.Name))
	err := no.Deployments().Delete(context.Background(), k.Name, v1.DeleteOptions{})
	if err != nil {
		logger.Error("unable to update deployment resource on kubernetes", zap.String("deployment-name", k.ObjectMeta.Name))
	}
	return err
}

func (no *NodeOperator) dappApplications(app *meta.App, usePermTree bool) []applyable {
	return []applyable{
		no.toSecret(app),
		no.dAppToDeployment(app, usePermTree),
		no.dappToService(app),
	}
}
