package add_k8s_labels_metadata

import (
	"fmt"
	"strings"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/processors"
)

const (
	processorName                          = "add_k8s_labels_metadata"
	keyKubernetesPodName                   = "kubernetes.pod.name"
	keyKubernetesNamespace                 = "kubernetes.namespace"
	keyKubernetesContainerName             = "kubernetes.container.name"
	keyKubernetesAnnotationsControllerKind = "kubernetes.annotations.k8s.labels.cloud/controller-kind"
	keyKubernetesLabelAppName              = "kubernetes.labels.app"
)

func init() {
	processors.RegisterPlugin("add_k8s_labels_metadata", newK8sLabelsProcessor)
}

type addK8sLabels struct {
	log *logp.Logger
}
//增加processes labels 对events数据改写。动态根据k8s 服务状态生成topic , 输出k8s deployment、daemonset、statfuleset、job、configmap、pod等对象数据
func newK8sLabelsProcessor(cfg *common.Config) (processors.Processor, error) {
	return &addK8sLabels{logp.NewLogger(processorName)}, nil
}

func (d *addK8sLabels) Run(event *beat.Event) (*beat.Event, error) {
	defaultControllerKind := "deployment"
	defaultTopicInfix := "docker"
	defaultDeploymentName := "none"
	defaultAppName := "none"

	kubernetesAnnotationsControllerKind, err := event.Fields.GetValue(keyKubernetesAnnotationsControllerKind)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s ,err %v", keyKubernetesAnnotationsControllerKind, event.Fields.String(),err)
	} else {
		defaultControllerKind = strings.ToLower(kubernetesAnnotationsControllerKind.(string))
	}

	if defaultControllerKind != "deployment" {
		defaultTopicInfix = defaultControllerKind
	}
	event.Fields["controller_kind"] = defaultControllerKind

	kubernetesLabelAppName, err := event.Fields.GetValue(keyKubernetesLabelAppName)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s ,err %v", keyKubernetesLabelAppName, event.Fields.String(),err)
	} else {
		defaultAppName = kubernetesLabelAppName.(string)
		if defaultControllerKind == "deployment" {
			defaultDeploymentName = defaultAppName
		}

	}

	// 为了兼容之前的deployment 和topic字段
	// 非 deployment 或者无法取到app label的deployment为none
	// topic 如果类型为deployment 则为 k8s_docker_{{deploymentName}}, 其他类型 则为 k8s_{{controllerKind}}_{{appName}}
	event.Fields["deployment"] = defaultDeploymentName
	event.Fields["topic"] = fmt.Sprintf("k8s_%s_%s", defaultTopicInfix, defaultAppName)
	event.Fields["app"] = defaultAppName

	kubernetesPodName, err := event.Fields.GetValue(keyKubernetesPodName)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s", keyKubernetesPodName, event.Fields.String())
	} else {
		event.Fields["pod"] = kubernetesPodName
	}

	kubernetesNamespace, err := event.Fields.GetValue(keyKubernetesNamespace)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s ", keyKubernetesNamespace, event.Fields.String())
	} else {
		event.Fields["namespace"] = kubernetesNamespace
	}

	kubernetesContainerName, err := event.Fields.GetValue(keyKubernetesContainerName)
	if err != nil {
		d.log.Debugf("Error while get %s fields. %s", keyKubernetesContainerName, event.Fields.String())
	} else {
		event.Fields["container_name"] = kubernetesContainerName
	}

	return event, nil
}

func (d *addK8sLabels) String() string {
	return processorName
}
