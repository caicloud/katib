package suggestion

import (
	suggestionv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createOrUpdateCondition(suggestionStatus *suggestionv1alpha2.SuggestionStatus,
	conditionType suggestionv1alpha2.SuggestionConditionType,
	conditionStatus corev1.ConditionStatus) {
	createOrUpdateConditionWithReason(suggestionStatus, conditionType, conditionStatus, "", "")
}

func createOrUpdateConditionWithReason(suggestionStatus *suggestionv1alpha2.SuggestionStatus,
	conditionType suggestionv1alpha2.SuggestionConditionType,
	conditionStatus corev1.ConditionStatus,
	reason, message string) {
	conditions := suggestionStatus.Conditions
	for i, cond := range conditions {
		if cond.Type == conditionType {
			updateCondition(&suggestionStatus.Conditions[i], conditionStatus, reason, message)
			return
		}
	}
	c := createCondition(conditionType, conditionStatus, reason, message)
	suggestionStatus.Conditions = append(suggestionStatus.Conditions, c)

}

func createCondition(conditionType suggestionv1alpha2.SuggestionConditionType,
	status corev1.ConditionStatus,
	reason string,
	message string) suggestionv1alpha2.SuggestionCondition {
	return suggestionv1alpha2.SuggestionCondition{
		Type:               conditionType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
	}
}

func updateCondition(condition *suggestionv1alpha2.SuggestionCondition,
	status corev1.ConditionStatus, reason string, message string) {
	if condition.Status != status {
		condition.LastTransitionTime = metav1.Now()
	}
	condition.Status = status
	condition.Reason = reason
	condition.Message = message
	condition.LastUpdateTime = metav1.Now()

}
