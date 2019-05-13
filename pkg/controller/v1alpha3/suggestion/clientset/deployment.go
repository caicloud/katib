package clientset

import (
	"context"
	"reflect"

	suggestionv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DeploymentClient struct {
	client.Client
}

func New(c client.Client) DeploymentClient {
	return DeploymentClient{
		Client: c,
	}
}

func (dc *DeploymentClient) CreateOrUpdateDeployment(suggestion *suggestionv1alpha2.Suggestion, desired *appsv1.Deployment) error {
	found := &appsv1.Deployment{}
	err := dc.Get(context.TODO(), types.NamespacedName{
		Name:      desired.Name,
		Namespace: desired.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// log
		// log.Info("Creating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
		err = dc.Create(context.TODO(), desired)
		return err
	} else if err != nil {
		return err
	}

	if !reflect.DeepEqual(desired.Spec, found.Spec) {
		// found.Spec = desired.Spec
		// log
		if err = dc.Update(context.TODO(), desired); err != nil {
			return err
		}
	}
	desired.Status = found.Status
	return nil
}
