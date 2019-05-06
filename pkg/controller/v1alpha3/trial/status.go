package trial

import (
	"context"

	commonv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/common/v1beta1"
	tfv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
)

func (r *ReconcileTrial) updateStatus(t *trialv1alpha2.Trial) error {
	return r.Status().Update(context.TODO(), t)
}

func (r *ReconcileTrial) updateStatusIfChanged(t *trialv1alpha2.Trial, oldStatus *trialv1alpha2.TrialStatus) error {
	if !equality.Semantic.DeepEqual(oldStatus, &t.Status) {
		return r.updateStatusHandler(t)
	}
	return nil
}

func (r *ReconcileTrial) SyncTFJobStatus(t *trialv1alpha2.Trial, tfJob *tfv1beta1.TFJob) error {
	jobStatus := tfJob.Status
	if len(jobStatus.Conditions) > 0 {
		lc := jobStatus.Conditions[len(jobStatus.Conditions)-1]
		if lc.Type == commonv1beta1.JobSucceeded &&
			lc.Status == corev1.ConditionTrue {
			createOrUpdateCondition(&t.Status, trialv1alpha2.TrialSucceeded, corev1.ConditionTrue)
			createOrUpdateCondition(&t.Status, trialv1alpha2.TrialRunning, corev1.ConditionFalse)
		} else if lc.Type == commonv1beta1.JobFailed &&
			lc.Status == corev1.ConditionTrue {
			createOrUpdateCondition(&t.Status, trialv1alpha2.TrialSucceeded, corev1.ConditionFalse)
			createOrUpdateCondition(&t.Status, trialv1alpha2.TrialFailed, corev1.ConditionTrue)
			createOrUpdateCondition(&t.Status, trialv1alpha2.TrialRunning, corev1.ConditionFalse)
		} else if lc.Type == commonv1beta1.JobRunning &&
			lc.Status == corev1.ConditionTrue {
			createOrUpdateCondition(&t.Status, trialv1alpha2.TrialRunning, corev1.ConditionTrue)
		}
	}
	return nil
}
