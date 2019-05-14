package suggestion

import (
	"fmt"

	suggestionv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/suggestion/v1alpha2"
	corev1 "k8s.io/api/core/v1"
)

func (r *ReconcileSuggestion) reportChange(s *suggestionv1alpha2.Suggestion, operator, typ string) {
	msg := fmt.Sprintf("%s the %s", operator, typ)
	log.Info(msg)
	r.Recorder.ReportChange(s, operator, typ)
}

func (r *ReconcileSuggestion) reportError(s *suggestionv1alpha2.Suggestion, err error, reason, msg string) {
	log.Error(err, msg)
	r.Recorder.ReportError(s, reason, err.Error())
	createOrUpdateConditionWithReason(&s.Status, suggestionv1alpha2.SuggestionDeploymentAvailable, corev1.ConditionFalse, reason, msg)
}
