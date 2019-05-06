package clientset

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/recorder"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/util"
)

const (
	loggerNameUnstructured = "unstructured-client"
)

// Unstructured is the type for unstructured client.
type Unstructured interface {
	CreateOrUpdateUnifiedJob(t *trialv1alpha2.Trial, job *unstructured.Unstructured) error
}

// GeneralUnstructured is the general client for Unstructured.
type GeneralUnstructured struct {
	client.Client
	recorder.Recorder
}

// CreateOrUpdateUnifiedJob creates or updates the unified job owned by the trial.
func (g *GeneralUnstructured) CreateOrUpdateUnifiedJob(t *trialv1alpha2.Trial, job *unstructured.Unstructured) error {
	typedName := types.NamespacedName{
		Name:      job.GetName(),
		Namespace: job.GetNamespace(),
	}
	logger := logf.Log.WithName(typedName.String())
	found := job.DeepCopy()
	err := g.Get(context.TODO(), typedName, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating Job", "namespace", job.GetNamespace(), "name", job.GetName())
		err = g.Create(context.TODO(), job)
		if err != nil {
			return err
		}
		g.ReportChange(t, util.FlagCreate, util.TypeTFJob)
		return nil
	} else if err != nil {
		return err
	}

	// We do not support updating now.
	return nil
}

// NewUnstructured creates a new Unstructured client.
func NewUnstructured(c client.Client, r recorder.Recorder) Unstructured {
	return &GeneralUnstructured{
		Client:   c,
		Recorder: r,
	}
}
