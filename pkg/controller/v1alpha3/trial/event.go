package trial

import (
	corev1 "k8s.io/api/core/v1"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
)

// reportError reports the error to the event, and set the status to unhealthy.
func (r *ReconcileTrial) reportError(
	t *trialv1alpha2.Trial, err error, msg, reason string) {
	log.Info(msg, "error", err)
	r.Event(t, corev1.EventTypeWarning, reason, err.Error())
	createOrUpdateConditionWithReason(&t.Status, trialv1alpha2.TrialSucceeded,
		corev1.ConditionFalse, reason, msg)
	createOrUpdateConditionWithReason(&t.Status, trialv1alpha2.TrialFailed,
		corev1.ConditionTrue, reason, msg)
}
