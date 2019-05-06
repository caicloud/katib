package trial

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
)

func hasCondition(status trialv1alpha2.TrialStatus,
	conditionType trialv1alpha2.TrialConditionType,
	boolVal corev1.ConditionStatus) bool {
	for _, c := range status.Conditions {
		if c.Type == conditionType && c.Status == boolVal {
			return true
		}
	}
	return false
}

func createOrUpdateConditionWithReason(status *trialv1alpha2.TrialStatus,
	conditionType trialv1alpha2.TrialConditionType,
	boolVal corev1.ConditionStatus, reason, msg string) {
	if !containConditionType(status, conditionType) {
		status.Conditions = append(status.Conditions, newCondition(conditionType,
			boolVal, reason, msg))
	} else {
		for i := range status.Conditions {
			if status.Conditions[i].Type == conditionType {
				if status.Conditions[i].Status != boolVal {
					status.Conditions[i].LastTransitionTime = metav1.Now()
				}
				status.Conditions[i].Status = boolVal
				status.Conditions[i].LastUpdateTime = metav1.Now()
				status.Conditions[i].Reason = reason
				status.Conditions[i].Message = msg
			}
		}
	}
}

func newCondition(conditionType trialv1alpha2.TrialConditionType,
	boolVal corev1.ConditionStatus,
	reason, message string) trialv1alpha2.TrialCondition {
	return trialv1alpha2.TrialCondition{
		Type:               conditionType,
		Status:             boolVal,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func createOrUpdateCondition(status *trialv1alpha2.TrialStatus,
	conditionType trialv1alpha2.TrialConditionType,
	boolVal corev1.ConditionStatus) {
	createOrUpdateConditionWithReason(status, conditionType, boolVal, "", "")
}

func containConditionType(status *trialv1alpha2.TrialStatus,
	conditionType trialv1alpha2.TrialConditionType) bool {
	for _, condition := range status.Conditions {
		if condition.Type == conditionType {
			return true
		}
	}
	return false
}
