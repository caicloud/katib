package trial

import (
	tfv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/util"
)

func (r *ReconcileTrial) handle(trial *trialv1alpha2.Trial) (result reconcile.Result, err error) {
	oldStatus := trial.Status.DeepCopy()
	logger := log.WithName(types.NamespacedName{
		Namespace: trial.Namespace,
		Name:      trial.Name,
	}.String())
	result = reconcile.Result{}

	defer func() {
		logger.V(0).Info("Updating the status", "status", trial.Status)
		err = r.updateStatusIfChanged(trial, oldStatus)
	}()

	typ, err := r.GetDesiredJobType(trial)
	if err != nil {
		r.reportError(trial, err, "Fail to compose the job specification", util.FailReason)
		return result, err
	}

	switch typ {
	case tfv1beta1.Kind:
		// If the trial is succeeded, update the final metrics.
		if hasCondition(trial.Status, trialv1alpha2.TrialSucceeded, corev1.ConditionTrue) {
			if err = r.CollectFinalMetric(trial); err != nil {
				r.reportError(trial, err, "Fail to collect the metrics", util.FailReason)
				return result, err
			}
			return result, err
		}

		job, err := r.GetDesiredTFJob(trial)
		if err != nil {
			r.reportError(trial, err, "Fail to compose the job specification", util.FailReason)
			return result, err
		}
		if err = controllerutil.SetControllerReference(trial, job, r.scheme); err != nil {
			r.reportError(trial, err, "Fail to set owner", util.FailReason)
			return result, err
		}
		if err = r.CreateOrUpdateTFJob(trial, job); err != nil {
			r.reportError(trial, err, "Fail to create the job", util.FailReason)
			return result, err
		}
		// TODO(gaocegege): Get the periodical metrics.
		// TODO(gaocegege): Communicate with early stopping services.
		if err = r.SyncTFJobStatus(trial, job); err != nil {
			r.reportError(trial, err, "Fail to sync the job status", util.FailReason)
			return result, err
		}
	}
	return result, err
}
