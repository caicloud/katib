package clientset

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"

	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/util"
)

type Trial interface {
	FilterRunningTrials(
		trials []trialv1alpha2.Trial) []trialv1alpha2.Trial
	FilterFailedTrials(
		trials []trialv1alpha2.Trial) []trialv1alpha2.Trial
	GetTrialsOwnedBy(e *experimentv1alpha2.Experiment) ([]trialv1alpha2.Trial, error)
}

type GeneralTrialClient struct {
	client.Client
}

// GetTrialsOwnedBy returns the trials owned by the experiment.
func (g *GeneralTrialClient) GetTrialsOwnedBy(
	e *experimentv1alpha2.Experiment) ([]trialv1alpha2.Trial, error) {
	actual := &trialv1alpha2.TrialList{}

	requirement, err := labels.NewRequirement(util.LabelExperiment,
		selection.Equals, []string{
			e.Name,
		})
	if err != nil {
		return nil, err
	}
	if err = g.List(context.TODO(), &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*requirement)}, actual); err != nil {
		return nil, err
	}

	return actual.Items, nil
}

// FilterRunningTrials returns the running trials owned by the experiment.
func (g *GeneralTrialClient) FilterRunningTrials(
	trials []trialv1alpha2.Trial) []trialv1alpha2.Trial {
	actual := make([]trialv1alpha2.Trial, 0)
	for _, t := range trials {
		if isTrialRunning(&t) {
			actual = append(actual, t)
		}
	}
	return actual
}

// FilterFailedTrials returns the failed trials owned by the experiment.
func (g *GeneralTrialClient) FilterFailedTrials(
	trials []trialv1alpha2.Trial) []trialv1alpha2.Trial {
	actual := make([]trialv1alpha2.Trial, 0)
	for _, t := range trials {
		if isTrialFailed(&t) {
			actual = append(actual, t)
		}
	}
	return actual
}

func isTrialRunning(t *trialv1alpha2.Trial) bool {
	for _, c := range t.Status.Conditions {
		if c.Type == trialv1alpha2.TrialRunning && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isTrialFailed(t *trialv1alpha2.Trial) bool {
	for _, c := range t.Status.Conditions {
		if c.Type == trialv1alpha2.TrialFailed && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// New returns a new Trial Client.
func New(c client.Client) Trial {
	return &GeneralTrialClient{
		Client: c,
	}
}
