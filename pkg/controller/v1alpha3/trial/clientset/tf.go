package clientset

import (
	"context"

	tfv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/recorder"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/util"
)

const (
	loggerNameTF = "tensorflow-client"
)

// TensorFlow is the type for tensorflow client.
type TensorFlow interface {
	CreateOrUpdateTFJob(t *trialv1alpha2.Trial, tfJob *tfv1beta1.TFJob) error
}

// GeneralTF is the general client for TensorFlow.
type GeneralTF struct {
	client.Client
	recorder.Recorder
}

// CreateOrUpdateTFJob creates or updates the TFJob owned by the trial.
func (g *GeneralTF) CreateOrUpdateTFJob(t *trialv1alpha2.Trial, tfJob *tfv1beta1.TFJob) error {
	log := logf.Log.WithName(loggerNameTF)
	found := &tfv1beta1.TFJob{}
	err := g.Get(context.TODO(), types.NamespacedName{
		Name:      tfJob.Name,
		Namespace: tfJob.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating TFJob", "namespace", tfJob.Namespace, "name", tfJob.Name)
		err = g.Create(context.TODO(), tfJob)
		if err != nil {
			return err
		}
		g.ReportChange(t, util.FlagCreate, util.TypeTFJob)
		return nil
	}

	// // Update the found object and write the result back if there are any changes.
	// if !reflect.DeepEqual(tfJob.Spec, found.Spec) {
	// 	found.Spec = tfJob.Spec
	// 	log.Info("Updating TFJob", "namespace", tfJob.Namespace, "name", tfJob.Name)
	// 	err = g.Update(context.TODO(), tfJob)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	g.ReportChange(t, util.FlagUpdate, util.TypeTFJob)
	// }
	tfJob.Status = found.Status
	return nil
}

// New creates a new TFJob client.
func NewTF(c client.Client, r recorder.Recorder) TensorFlow {
	return &GeneralTF{
		Client:   c,
		Recorder: r,
	}
}
