package suggestion

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"

	suggestionsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
	suggestionv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
)

func (r *ReconcileSuggestion) syncStatus(deployStatus *appsv1.DeploymentStatus, suggestion *suggestionv1alpha2.Suggestion) error {
	for _, cond := range deployStatus.Conditions {
		if cond.Type == appsv1.DeploymentAvailable {
			createOrUpdateCondition(&suggestion.Status, suggestionv1alpha2.SuggestionDeploymentAvailable, cond.Status)
		} else if cond.Type == appsv1.DeploymentProgressing {
			createOrUpdateCondition(&suggestion.Status, suggestionv1alpha2.SuggestionDeploymentProgressing, cond.Status)
		} else if cond.Type == appsv1.DeploymentReplicaFailure {
			createOrUpdateCondition(&suggestion.Status, suggestionv1alpha2.SuggestionDeploymentReplicaFailure, cond.Status)
		}
	}
	return nil
}

func (r *ReconcileSuggestion) updateStatusIfChanged(oldStatus *suggestionv1alpha2.SuggestionStatus,
	suggestion *suggestionsv1alpha2.Suggestion) error {
	if !equality.Semantic.DeepEqual(oldStatus, &suggestion.Status) {
		return r.Status().Update(context.TODO(), suggestion)
	}
	return nil

}
