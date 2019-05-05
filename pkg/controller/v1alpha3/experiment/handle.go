package experiment

import (
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
)

func (r *ReconcileExperiment) handle(experiment *experimentv1alpha2.Experiment) (result reconcile.Result, err error) {
	oldStatus := experiment.Status.DeepCopy()
	result = reconcile.Result{}
	defer func() {
		err = r.updateStatusIfChanged(experiment, oldStatus)
	}()

	trials, err := r.GetTrialsOwnedBy(experiment)
	if err != nil {
		return result, err
	}
	failedTrials := r.FilterFailedTrials(trials)
	runningTrials := r.FilterRunningTrials(trials)
	lenTrials := len(trials)
	lenFailedTrials := len(failedTrials)
	lenRunningTrials := len(runningTrials)
	// TODO(gaocegege): Set all statistics.
	r.setStatistics(experiment, lenTrials, lenFailedTrials, lenRunningTrials)

	// The job is completed, we mark it and return.
	if experiment.Spec.MaxTrialCount != nil {
		if lenTrials >= *experiment.Spec.MaxTrialCount {
			r.markCompleted(experiment)
			return result, err
		}
	}

	if experiment.Spec.MaxFailedTrialCount != nil {
		if lenFailedTrials >= *experiment.Spec.MaxFailedTrialCount {
			r.markFailed(experiment)
			return result, err
		}
	}

	// Reconcile trials.
	if experiment.Spec.ParallelTrialCount != nil {
		if len(runningTrials) < *experiment.Spec.ParallelTrialCount {
			// TODO(gaocegege): Create new trials.
		}
	}

	// TODO(gaocegege): Update best trial.

	return result, err
}
