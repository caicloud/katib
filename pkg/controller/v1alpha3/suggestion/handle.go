package suggestion

import (
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/util"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileSuggestion) handle(suggestion *suggestionsv1alpha2.Suggestion) (result reconcile.Result, err error) {
	oldStatus := suggestion.Status.DeepCopy()
	result = reconcile.Result{}
	logger := log.WithName(types.NamespacedName{
		Namespace: suggestion.Namespace,
		Name:      suggestion.Name,
	}.String())

	defer func() {
		err = r.updateStatusIfChanged(oldStatus, suggestion)
	}()

	desired, err := getDesiredDeployment(suggestion)
	if err != nil {
		r.reportError(suggestion, err, util.FailReason, "Failed to get desired deployment")
		return result, err
	}
	logger.V(0).Info("OK: Get desired deployment of suggestion")

	if err = controllerutil.SetControllerReference(suggestion, desired, r.scheme); err != nil {
		r.reportError(suggestion, err, util.FailReason, "Failed to set controller reference")
		return result, err
	}
	logger.V(0).Info("OK: Set controller reference between deployment and suggestion")

	// if suggestion spec changes, create or update deployment
	// desired deployment status is updated
	if err = r.CreateOrUpdateDeployment(suggestion, desired); err != nil {
		r.reportError(suggestion, err, util.FailReason, "Failed to create or update deployment of suggestion")
		return result, err
	}
	logger.V(0).Info("OK: Create or update deployment")

	// if deployment changes, sync status of suggestion
	if err = r.syncStatus(&desired.Status, suggestion); err != nil {
		r.reportError(suggestion, err, util.FailReason, "Failed to sync status of suggestion")
		return result, err
	}
	logger.V(0).Info("OK: Sync status of suggestion")

	return result, err
}
