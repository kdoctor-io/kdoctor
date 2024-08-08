// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package k8sObjManager

import (
	"context"
	"errors"

	"github.com/kdoctor-io/kdoctor/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetPodList(ctx context.Context, opts ...client.ListOption) ([]corev1.Pod, error) {
	var podlist corev1.PodList
	if e := nm.client.List(ctx, &podlist, opts...); e != nil {
		return nil, e
	}
	return podlist.Items, nil
}

func (nm *k8sObjManager) ListSelectedPod(ctx context.Context, labelSelector *metav1.LabelSelector) ([]corev1.Pod, error) {

	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, err
	}

	return nm.GetPodList(
		ctx,
		client.MatchingLabelsSelector{Selector: selector},
	)
}

func (nm *k8sObjManager) ListSelectedPodIPs(ctx context.Context, labelSelector *metav1.LabelSelector) (PodIps, error) {

	podlist, e := nm.ListSelectedPod(ctx, labelSelector)
	if e != nil {
		return nil, e
	}
	if len(podlist) == 0 {
		return nil, errors.New("failed to get any pods")
	}

	result := PodIps{}

	for _, v := range podlist {
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

func (nm *k8sObjManager) ListSelectedPodMultusIPs(ctx context.Context, labelSelector *metav1.LabelSelector) (PodIps, error) {
	podlist, e := nm.ListSelectedPod(ctx, labelSelector)
	if e != nil {
		return nil, e
	}
	if len(podlist) == 0 {
		return nil, errors.New("failed to get any pods")
	}

	return parseMultusIP(podlist)
}
