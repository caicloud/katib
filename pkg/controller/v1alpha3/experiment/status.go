package experiment

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"

	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
)

func (r *ReconcileExperiment) updateStatus(e *experimentv1alpha2.Experiment) error {
	return r.Status().Update(context.TODO(), e)
}

func (r *ReconcileExperiment) updateStatusIfChanged(e *experimentv1alpha2.Experiment, oldStatus *experimentv1alpha2.ExperimentStatus) error {
	if !equality.Semantic.DeepEqual(oldStatus, &e.Status) {
		return r.updateStatusHandler(e)
	}
	return nil
}

func (r *ReconcileExperiment) setStatistics(e *experimentv1alpha2.Experiment,
	trials, failedTrails, runningTrials int) {
	e.Status.TrialsFailed = failedTrails
	e.Status.TrialsRunning = runningTrials
	e.Status.Trials = trials
}

func (r *ReconcileExperiment) markCompleted(e *experimentv1alpha2.Experiment) {
	createOrUpdateCondition(&e.Status, experimentv1alpha2.ExperimentCompleted, corev1.ConditionTrue)
}

func (r *ReconcileExperiment) markFailed(e *experimentv1alpha2.Experiment) {
	createOrUpdateCondition(&e.Status, experimentv1alpha2.ExperimentFailed, corev1.ConditionTrue)
	createOrUpdateCondition(&e.Status, experimentv1alpha2.ExperimentCompleted, corev1.ConditionFalse)
}
