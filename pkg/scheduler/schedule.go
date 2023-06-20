package scheduler

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
)

/*
pod:
- name: ${task}-agent
- ns: controller-ns
- annotation:
- labels
- hostNetwork
- restartPolicy
- tolerations
- nodeSelector
- volumes
container:
- name:
- ns:
- image:
- imagePullPolicy:
- command:
- args:
- ports:
- startupProbe,livenessProbe,readinessProbe
- resources:
- env:
- volumeMounts:
*/

const (
	containerImage         = "ghcr.io/kdoctor-io/kdoctor-agent"
	containerCommand       = "/usr/bin/agent"
	containerArgConfigPath = "--config-path=/tmp/config-map/conf.yml"
	containerArgTaskKind   = "--task-kind"
	containerArgTaskName   = "--task-name"
	imagePullPolicy        = corev1.PullIfNotPresent
	configmapVolumeName    = "config-path"
	reportVolumeName       = "report-data"
	configmapName          = "kdoctor"
)

var (
	// container port
	containerPort = corev1.ContainerPort{
		Name:          "http",
		ContainerPort: 80,
		Protocol:      corev1.ProtocolTCP,
	}

	startupProbe = corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/healthy/startup",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(80),
				},
				Scheme: corev1.URISchemeHTTP,
			},
		},
		PeriodSeconds:    2,
		SuccessThreshold: 1,
		FailureThreshold: 60,
	}
	livenessProbe = corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/healthy/liveness",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(80),
				},
				Scheme: corev1.URISchemeHTTP,
			},
		},
		InitialDelaySeconds: 60,
		TimeoutSeconds:      5,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    6,
	}
	readinessProbe = corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/healthy/readiness",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: int32(80),
				},
				Scheme: corev1.URISchemeHTTP,
			},
		},
		InitialDelaySeconds: 5,
		TimeoutSeconds:      5,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}

	configmapVolumeMount = corev1.VolumeMount{
		Name:      configmapVolumeName,
		ReadOnly:  true,
		MountPath: "/tmp/config-map",
	}

	reportVolumeMount = corev1.VolumeMount{
		Name:      reportVolumeName,
		ReadOnly:  false,
		MountPath: "/report",
	}

	configmapVolume = corev1.Volume{
		Name: configmapVolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configmapName,
				},
				DefaultMode: pointer.Int32(0400),
			},
		},
	}

	hostPathDOC      = corev1.HostPathDirectoryOrCreate
	reportHostVolume = corev1.Volume{
		Name: reportVolumeName,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/run/kdoctor/agent",
				Type: &hostPathDOC,
			},
		},
	}

	podLabels = map[string]string{
		"app.kubernetes.io/component": "kdoctor-agent",
	}
)

func CreateTaskRuntimeIfNotExist(ctx context.Context, clientSet client.Client, taskKind, appNS string, task metav1.Object, agentSpec v1beta1.AgentSpec, log *zap.Logger) error {
	appName := fmt.Sprintf("kdoctor-%s", task.GetName())

	var app client.Object
	if agentSpec.Kind == "Deployment" {
		app = &appsv1.Deployment{}
	} else {
		app = &appsv1.DaemonSet{}
	}

	var needCreate bool
	objectKey := client.ObjectKey{
		Namespace: appNS,
		Name:      appName,
	}

	log.Sugar().Infof("try to get task %s cooresponding runtime %s/%s", task.GetName(), agentSpec.Kind, appName)
	err := clientSet.Get(ctx, objectKey, app)
	if nil != err {
		if errors.IsNotFound(err) {
			log.Sugar().Infof("task %s cooresponding runtime %s/%s not found, try to create one", task.GetName(), agentSpec.Kind, appName)
			needCreate = true
		} else {
			return err
		}
	}

	if needCreate {
		if agentSpec.Kind == "Deployment" {
			app = generateDeployment(appNS, appName, taskKind, task.GetName(), agentSpec)
		} else {
			app = generateDaemonSet(appNS, appName, taskKind, task.GetName(), agentSpec)
		}

		err := controllerruntime.SetControllerReference(task, app, clientSet.Scheme())
		if nil != err {
			return fmt.Errorf("failed to set controllerReference for %s %+v", agentSpec.Kind, app)
		}

		log.Sugar().Infof("try to create %s %+v for task %s", agentSpec.Kind, app, task.GetName())
		err = clientSet.Create(ctx, app)
		if nil != err {
			return err
		}
	}

	return nil
}

func DeleteTaskRuntime(ctx context.Context, clientSet client.Client, taskNS, taskName, kind string, log *zap.Logger) error {
	appName := fmt.Sprintf("kdoctor-%s", taskName)

	var app client.Object
	if kind == "Deployment" {
		app = &appsv1.Deployment{}
	} else {
		app = &appsv1.DaemonSet{}
	}
	app.SetNamespace(taskNS)
	app.SetName(appName)

	err := clientSet.Delete(ctx, app)
	return client.IgnoreNotFound(err)
}

func generateOneContainer(containerName, taskKind, taskName string, agentSpec v1beta1.AgentSpec) corev1.Container {
	var resource corev1.ResourceRequirements
	if agentSpec.Resources != nil {
		resource = *agentSpec.Resources
	}

	container := corev1.Container{
		Name:    containerName,
		Image:   containerImage,
		Command: []string{containerCommand},
		Args: []string{
			containerArgConfigPath,
			fmt.Sprintf(containerArgTaskKind+"=%s", taskKind),
			fmt.Sprintf(containerArgTaskName+"=%s", taskName),
		},
		Ports:           []corev1.ContainerPort{containerPort},
		Env:             agentSpec.Env,
		Resources:       resource,
		VolumeMounts:    []corev1.VolumeMount{configmapVolumeMount, reportVolumeMount},
		LivenessProbe:   &livenessProbe,
		ReadinessProbe:  &readinessProbe,
		StartupProbe:    &startupProbe,
		ImagePullPolicy: imagePullPolicy,
	}

	return container
}

func generateDaemonSet(appNS, appName, taskKind, taskName string, agentSpec v1beta1.AgentSpec) *appsv1.DaemonSet {
	container := generateOneContainer(appName, taskKind, taskName, agentSpec)

	var volumes []corev1.Volume
	volumes = append(volumes, configmapVolume)
	if agentSpec.PersistentVolumeClaim != nil {
		pvcVolume := corev1.Volume{
			Name: reportVolumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: agentSpec.PersistentVolumeClaim,
			},
		}
		volumes = append(volumes, pvcVolume)
	} else {
		volumes = append(volumes, reportHostVolume)
	}

	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        appName,
			Namespace:   appNS,
			Labels:      nil,
			Annotations: agentSpec.Annotation,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: podLabels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      podLabels,
					Annotations: agentSpec.Annotation,
				},
				Spec: corev1.PodSpec{
					Volumes:                       volumes,
					Containers:                    []corev1.Container{container},
					RestartPolicy:                 corev1.RestartPolicyAlways,
					TerminationGracePeriodSeconds: agentSpec.TerminationGracePeriodSeconds,
					HostNetwork:                   agentSpec.HostNetwork,
					Affinity:                      nil,
				},
			},
		},
	}

	return daemonSet
}

func generateDeployment(appNS, appName, taskKind, taskName string, agentSpec v1beta1.AgentSpec) *appsv1.Deployment {
	container := generateOneContainer(appName, taskKind, taskName, agentSpec)

	var volumes []corev1.Volume
	volumes = append(volumes, configmapVolume)
	if agentSpec.PersistentVolumeClaim != nil {
		pvcVolume := corev1.Volume{
			Name: reportVolumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: agentSpec.PersistentVolumeClaim,
			},
		}
		volumes = append(volumes, pvcVolume)
	} else {
		volumes = append(volumes, reportHostVolume)
	}

	deployAffinity := &corev1.Affinity{
		NodeAffinity: nil,
		PodAffinity:  nil,
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: nil,
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				corev1.WeightedPodAffinityTerm{
					Weight: 100,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchLabels:      podLabels,
							MatchExpressions: nil,
						},
						Namespaces:        nil,
						TopologyKey:       "kubernetes.io/hostname",
						NamespaceSelector: nil,
					},
				},
			},
		},
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        appName,
			Namespace:   appNS,
			Labels:      nil,
			Annotations: agentSpec.Annotation,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: agentSpec.DeploymentReplicas,
			Selector: &metav1.LabelSelector{MatchLabels: podLabels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: agentSpec.Annotation,
					Labels:      podLabels,
				},
				Spec: corev1.PodSpec{
					Volumes:                       volumes,
					Containers:                    []corev1.Container{container},
					RestartPolicy:                 corev1.RestartPolicyAlways,
					TerminationGracePeriodSeconds: agentSpec.TerminationGracePeriodSeconds,
					HostNetwork:                   agentSpec.HostNetwork,
					Affinity:                      deployAffinity,
				},
			},
		},
	}

	return deployment
}
