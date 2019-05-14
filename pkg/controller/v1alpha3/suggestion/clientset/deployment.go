package clientset

import (
	"context"

	suggestionv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/recorder"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/util"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	clientsetName = "suggestion-deployment-clientset"
)

var log = logf.Log.WithName(clientsetName)

type DeploymentClient struct {
	client.Client
	recorder.Recorder
}

func New(c client.Client, r recorder.Recorder) DeploymentClient {
	return DeploymentClient{
		Client:   c,
		Recorder: r,
	}
}

func (dc *DeploymentClient) CreateOrUpdateDeployment(suggestion *suggestionv1alpha2.Suggestion, desired *appsv1.Deployment) error {
	found := &appsv1.Deployment{}
	err := dc.Get(context.TODO(), types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Deployment", "namespace", desired.Namespace, "name", desired.Name)
		if err = dc.Create(context.TODO(), desired); err != nil {
			return err
		}
		dc.ReportChange(suggestion, util.FlagCreate, util.TypeDeployment)
	} else if err != nil {
		return err
	}

	// TODO(anchovYu): check and update
	// if !reflect.DeepEqual(desired.Spec, found.Spec) {
	// 	// found.Spec = desired.Spec
	// 	log.Info("Updating Deployment", "namespace", desired.Namespace, "name", desired.Name)
	// 	if err = dc.Update(context.TODO(), desired); err != nil {
	// 		return err
	// 	}
	// 	dc.ReportChange(suggestion, util.FlagUpdate, util.TypeDeployment)
	// }
	desired.Status = found.Status
	return nil
}
