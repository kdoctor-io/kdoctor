// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package k8sObjManager

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetDaemonset(ctx context.Context, name, namespace string) (*appsv1.DaemonSet, error) {
	d := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	key := client.ObjectKeyFromObject(d)
	if e := nm.client.Get(ctx, key, d); e != nil {
		return nil, fmt.Errorf("failed to get daemonset %v/%v, reason=%v", namespace, name, e)
	}
	return d, nil
}

func (nm *k8sObjManager) ListDaemonsetPod(ctx context.Context, daemonsetName, daemonsetNameSpace string) ([]corev1.Pod, error) {

	dae, e := nm.GetDaemonset(ctx, daemonsetName, daemonsetNameSpace)
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

func (nm *k8sObjManager) ListDaemonsetPodNodes(ctx context.Context, daemonsetName, daemonsetNameSpace string) ([]string, error) {

	podlist, e := nm.ListDaemonsetPod(ctx, daemonsetName, daemonsetNameSpace)
	if e != nil {
		return nil, e
	}

	nodelist := []string{}
	for _, v := range podlist {
		nodelist = append(nodelist, v.Spec.NodeName)
	}
	return nodelist, nil
}

// --------------

func parseMultusIP(podlist []corev1.Pod) (PodIps, error) {

	result := PodIps{}

	MultusPodAnnotationKey := types.AgentConfig.Configmap.MultusPodAnnotationKey
	if len(types.ControllerConfig.Configmap.MultusPodAnnotationKey) > 0 {
		MultusPodAnnotationKey = types.ControllerConfig.Configmap.MultusPodAnnotationKey
	}

	// for eth0
	for _, v := range podlist {
		t := IPs{}
		t.InterfaceName = "eth0"
		for _, m := range v.Status.PodIPs {
			if utils.CheckIPv4Format(m.IP) {
				t.IPv4 = m.IP
			} else {
				t.IPv6 = m.IP
			}
		}
		result[v.Name] = []IPs{t}
	}

	// for other interface
	for _, v := range podlist {
		val, ok := v.Annotations[MultusPodAnnotationKey]
		if !ok {
			continue
		}
		tmp := []MultusAnnotationValueItem{}
		if err := json.Unmarshal([]byte(val), &tmp); err != nil {
			return nil, fmt.Errorf("failed to parse multus annotation, '%v',  error=%v", val, err)
		}
		for _, r := range tmp {
			if r.Interface == "eth0" || len(r.Ips) == 0 {
				continue
			}
			t := IPs{}
			t.InterfaceName = r.Interface
			for _, w := range r.Ips {
				if utils.CheckIPv4Format(w) {
					t.IPv4 = w
				} else {
					t.IPv6 = w
				}
			}
			result[v.Name] = append(result[v.Name], t)
		}

	}
	return result, nil
}

func (nm *k8sObjManager) ListDaemonsetPodMultusIPs(ctx context.Context, daemonsetName, daemonsetNameSpace string) (PodIps, error) {
	podlist, e := nm.ListDaemonsetPod(ctx, daemonsetName, daemonsetNameSpace)
	if e != nil {
		return nil, e
	}
	if len(podlist) == 0 {
		return nil, fmt.Errorf("failed to get any pods")
	}
	return parseMultusIP(podlist)
}

func (nm *k8sObjManager) ListDaemonsetPodIPs(ctx context.Context, daemonsetName, daemonsetNameSpace string) (PodIps, error) {

	podlist, e := nm.ListDaemonsetPod(ctx, daemonsetName, daemonsetNameSpace)
	if e != nil {
		return nil, e
	}
	if len(podlist) == 0 {
		return nil, fmt.Errorf("failed to get any pods")
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
