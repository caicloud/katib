package suggestion

import (
	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileSuggestion) handle(instance *suggestionsv1alpha2.Suggestion) (reconcile.Result, error) {

	return reconcile.Result{}, nil
}
