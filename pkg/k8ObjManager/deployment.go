// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package k8sObjManager

import (
	"context"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetDeployment(ctx context.Context, name, namespace string) (*appsv1.Deployment, error) {
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	key := client.ObjectKeyFromObject(d)
	if e := nm.client.Get(ctx, key, d); e != nil {
		return nil, fmt.Errorf("failed to get deployment %v/%v, reason=%v", namespace, name, e)
	}
	return d, nil
}

func (nm *k8sObjManager) ListDeploymentPod(ctx context.Context, deploymentName, deploymentNameSpace string) ([]corev1.Pod, error) {

	dae, e := nm.GetDeployment(ctx, deploymentName, deploymentNameSpace)
	if e != nil {
		return nil, fmt.Errorf("failed to get daemonset, error=%v", e)
	}

	podLable := dae.Spec.Template.Labels
	opts := []client.ListOption{
		client.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(podLable),
		},
	}
	return nm.GetPodList(ctx, opts...)
}

func (nm *k8sObjManager) ListDeploymentPodIPs(ctx context.Context, deploymentName, deploymentNameSpace string) (PodIps, error) {

	podList, e := nm.ListDeploymentPod(ctx, deploymentName, deploymentNameSpace)
	if e != nil {
		return nil, e
	}
	if len(podList) == 0 {
		return nil, fmt.Errorf("failed to get any pods")
	}

	result := PodIps{}

	for _, v := range podList {
		t := IPs{}
		t.InterfaceName = "eth0"
		for _, m := range v.Status.PodIPs {
			if utils.CheckIPv4Format(m.IP) {
				t.IPv4 = m.IP
			} else if utils.CheckIPv6Format(m.IP) {
				t.IPv6 = m.IP
			}
		}
		result[v.Name] = []IPs{t}
	}
	return result, nil
}

func (nm *k8sObjManager) ListDeployPodMultusIPs(ctx context.Context, deploymentName, deploymentNameSpace string) (PodIps, error) {
	podlist, e := nm.ListDeploymentPod(ctx, deploymentName, deploymentNameSpace)
	if e != nil {
		return nil, e
	}
	if len(podlist) == 0 {
		return nil, fmt.Errorf("failed to get any pods")
	}
	return parseMultusIP(podlist)
}
