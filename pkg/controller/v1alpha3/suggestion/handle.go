package suggestion

import (
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileSuggestion) handle(suggestion *suggestionsv1alpha2.Suggestion) (result reconcile.Result, err error) {
	oldStatus := suggestion.Status.DeepCopy()
	result = reconcile.Result{}

	defer func() {
		err = r.updateStatusIfChanged(oldStatus, suggestion)
	}()

	desired, err := getDesiredDeployment(suggestion)
	if err != nil {
		// TODO: log
		return result, err
	}

	if err = controllerutil.SetControllerReference(suggestion, desired, r.scheme); err != nil {
		// TODO: log
		return result, err
	}

	// if suggestion spec changes, create or update deployment
	// desired deployment status is updated
	if err = r.CreateOrUpdateDeployment(suggestion, desired); err != nil {
		// TODO: log
		return result, err
	}

	// if deployment changes, sync status of suggestion
	if err = r.syncStatus(&desired.Status, suggestion); err != nil {
		// TODO: log
		return result, err
	}

	return result, err
}
