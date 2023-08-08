// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package k8sObjManager

import (
	"context"
	"errors"
	"github.com/kdoctor-io/kdoctor/pkg/lock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sObjManager interface {
	// node
	GetNode(ctx context.Context, nodeName string) (*corev1.Node, error)
	ListNodes(ctx context.Context, opts ...client.ListOption) (*corev1.NodeList, error)
	MatchNodeSelected(ctx context.Context, nodeName string, labelSelector *metav1.LabelSelector) (bool, error)
	ListSelectedNodes(ctx context.Context, labelSelector *metav1.LabelSelector) ([]string, error)
	GetNodeIP(ctx context.Context, nodeName string) (ipv4, ipv6 string, err error)

	// daemonset
	ListDaemonsetPodNodes(ctx context.Context, daemonsetName, daemonsetNameSpace string) ([]string, error)
	GetDaemonset(ctx context.Context, name, namespace string) (*appsv1.DaemonSet, error)
	ListDaemonsetPodIPs(ctx context.Context, daemonsetName, daemonsetNameSpace string) (PodIps, error)
	ListDaemonsetPodMultusIPs(ctx context.Context, daemonsetName, daemonsetNameSpace string) (PodIps, error)

	// deployment
	GetDeployment(ctx context.Context, name, namespace string) (*appsv1.Deployment, error)
	ListDeploymentPodIPs(ctx context.Context, deploymentName, deploymentNameSpace string) (PodIps, error)
	ListDeployPodMultusIPs(ctx context.Context, deploymentName, deploymentNameSpace string) (PodIps, error)

	// pod
	GetPodList(ctx context.Context, opts ...client.ListOption) ([]corev1.Pod, error)
	ListSelectedPodMultusIPs(ctx context.Context, labelSelector *metav1.LabelSelector) (PodIps, error)
	ListSelectedPodIPs(ctx context.Context, labelSelector *metav1.LabelSelector) (PodIps, error)
	ListSelectedPod(ctx context.Context, labelSelector *metav1.LabelSelector) ([]corev1.Pod, error)

	// service
	GetService(ctx context.Context, name, namespace string) (*corev1.Service, error)
	GetServiceAccessUrl(ctx context.Context, name, namespace string, portName string) (*ServiceAccessUrl, error)
	ListServicesDnsIP(ctx context.Context) ([]string, error)

	GetIngress(ctx context.Context, name, namespace string) (*networkingv1.Ingress, error)

	// secret
	GetSecret(ctx context.Context, name, namespace string) (*corev1.Secret, error)

	// configmap
	GetConfigMap(ctx context.Context, name, namespace string) (*corev1.ConfigMap, error)
}

type k8sObjManager struct {
	client client.Client
}

var l lock.Mutex
var globalManager *k8sObjManager

func Initk8sObjManager(client client.Client) error {
	if client == nil {
		return errors.New("k8s client must be specified")
	}
	l.Lock()
	defer l.Unlock()

	if globalManager == nil {
		globalManager = &k8sObjManager{
			client: client,
		}
	}
	return nil
}

func GetK8sObjManager() K8sObjManager {
	if globalManager == nil {
		panic("globalManager is not initialize")
	}
	return globalManager
}
