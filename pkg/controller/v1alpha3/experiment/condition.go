package experiment

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
)

func createOrUpdateConditionWithReason(status *experimentv1alpha2.ExperimentStatus,
	conditionType experimentv1alpha2.ExperimentConditionType,
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

func newCondition(conditionType experimentv1alpha2.ExperimentConditionType,
	boolVal corev1.ConditionStatus,
	reason, message string) experimentv1alpha2.ExperimentCondition {
	return experimentv1alpha2.ExperimentCondition{
		Type:               conditionType,
		Status:             boolVal,
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func createOrUpdateCondition(status *experimentv1alpha2.ExperimentStatus,
	conditionType experimentv1alpha2.ExperimentConditionType,
	boolVal corev1.ConditionStatus) {
	createOrUpdateConditionWithReason(status, conditionType, boolVal, "", "")
}

func containConditionType(status *experimentv1alpha2.ExperimentStatus,
	conditionType experimentv1alpha2.ExperimentConditionType) bool {
	for _, condition := range status.Conditions {
		if condition.Type == conditionType {
			return true
		}
	}
	return false
}
