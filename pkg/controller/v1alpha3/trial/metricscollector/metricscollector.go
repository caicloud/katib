package metricscollector

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/util"
)

const (
	RecommendedKubeConfigPathEnv = "KUBECONFIG"
)

var (
	log             = logf.Log.WithName("metrics-collector")
	splitChar       = []string{" ", "=", " = ", ":", ": "}
	tailLines int64 = 1000
)

type MetricsCollector interface {
	CollectFinalMetric(trial *trialv1alpha2.Trial) error
}

type GeneralMetricsCollector struct {
	clientset *kubernetes.Clientset
}

func New() (MetricsCollector, error) {
	var kubeconfig string
	// Note: ENV KUBECONFIG will overwrite user defined Kubeconfig option.
	if len(os.Getenv(RecommendedKubeConfigPathEnv)) > 0 {
		// use the current context in kubeconfig
		// This is very useful for running locally.
		kubeconfig = os.Getenv(RecommendedKubeConfigPathEnv)
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &GeneralMetricsCollector{
		clientset: clientset,
	}, nil

}

func (d *GeneralMetricsCollector) CollectFinalMetric(trial *trialv1alpha2.Trial) error {
	// Initialize the observation.
	if trial.Status.Observation == nil {
		trial.Status.Observation = &trialv1alpha2.Observation{
			Metrics: make([]trialv1alpha2.Metric, 0),
		}
	}

	// Compose the query labelMap.
	labelMap := make(map[string]string)
	labelMap[util.LabelTrial] = trial.Name

	pl, err := d.clientset.CoreV1().Pods(trial.Namespace).List(metav1.ListOptions{LabelSelector: labels.Set(labelMap).String()})
	if err != nil {
		return err
	}
	if len(pl.Items) == 0 {
		return fmt.Errorf("No Pods are found in Trial %v", trial.Name)
	}
	logs, err := d.clientset.CoreV1().Pods(trial.Namespace).GetLogs(
		pl.Items[0].Name, &apiv1.PodLogOptions{TailLines: &tailLines}).Do().Raw()
	if err != nil {
		return err
	}
	if len(logs) == 0 {
		return fmt.Errorf("No logs are found in Trial %v", trial.Name)
	}
	return d.parseLogs(
		trial, strings.Split(string(logs), "\n"), trial.Spec.Objective, trial.Spec.Metrics)
}

func (d *GeneralMetricsCollector) parseLogs(trial *trialv1alpha2.Trial, logs []string, objectiveValueName string, metrics []string) error {
	logger := log.WithName(types.NamespacedName{
		Namespace: trial.Namespace,
		Name:      trial.Name,
	}.String())

	for _, logline := range logs {
		value, ok := canGet(logline, objectiveValueName)
		if ok {
			logger.V(0).Info("Extract metrics",
				"metrics", objectiveValueName, "value", value)
			trial.Status.Observation.Objective = &trialv1alpha2.Metric{
				Name:  objectiveValueName,
				Value: value,
			}
			continue
		}

		for _, m := range metrics {
			value, ok = canGet(logline, m)
			if ok {
				logger.V(0).Info("Extract metrics", "metrics", m, "value", value)
				setMetric(trial.Status.Observation, m, value)
				break
			}
		}
	}
	return nil
}

func canGet(line, key string) (float64, bool) {
	if strings.Contains(line, key) {
		for _, s := range splitChar {
			strs := strings.Split(line, s)
			if len(strs) != 2 {
				continue
			}
			if strs[0] != key {
				return 0, false
			}
			n, err := strconv.ParseFloat(strs[1], 64)
			if err != nil {
				return 0, false
			}
			return n, true
		}
	}
	return 0, false
}

func setMetric(observation *trialv1alpha2.Observation, name string, value float64) {
	for i, m := range observation.Metrics {
		if m.Name == name {
			observation.Metrics[i].Value = value
			return
		}
	}
	observation.Metrics = append(
		observation.Metrics, trialv1alpha2.Metric{
			Name:  name,
			Value: value,
		})
}
