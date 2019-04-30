package recorder

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// Recorder is the interface used to record events.
type Recorder interface {
	record.EventRecorder
	ReportChange(o runtime.Object, operator, typ string)
}

// GeneralRecorder is the general-purpose implementation of Recorder.
type GeneralRecorder struct {
	record.EventRecorder
}

// ReportChange reports the change to the object.
func (g *GeneralRecorder) ReportChange(o runtime.Object, operator, typ string) {
	msg := fmt.Sprintf("%s the %s", operator, typ)
	g.EventRecorder.Event(o, corev1.EventTypeNormal, "SuccessfulChange", msg)
}

// New returns a new Recorder.
func New(r record.EventRecorder) Recorder {
	return &GeneralRecorder{
		EventRecorder: r,
	}
}
