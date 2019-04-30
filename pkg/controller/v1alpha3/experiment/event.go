package experiment

import (
	corev1 "k8s.io/api/core/v1"

	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
)

// reportError reports the error to the event, and set the status to unhealthy.
func (r *ReconcileExperiment) reportError(
	e *experimentv1alpha2.Experiment, err error, msg, reason string) {
	log.Info(msg, "error", err)
	r.Event(e, corev1.EventTypeWarning, reason, err.Error())
	createOrUpdateConditionWithReason(&e.Status, experimentv1alpha2.ExperimentAvailable,
		corev1.ConditionFalse, reason, msg)
}
