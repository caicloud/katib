package composer

import (
	"bytes"

	tfv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/util"
)

var log = logf.Log.WithName("composer")

type Composer interface {
	GetDesiredJobType(instance *trialv1alpha2.Trial) (string, error)
	GetDesiredTFJob(instance *trialv1alpha2.Trial) (*tfv1beta1.TFJob, error)
}

type GeneralComposer struct {
}

func (g *GeneralComposer) GetDesiredTFJob(instance *trialv1alpha2.Trial) (*tfv1beta1.TFJob, error) {
	bufSize := 1024
	logger := log.WithValues("trial", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	buf := bytes.NewBufferString(instance.Spec.RunSpec)

	desiredJobSpec := &tfv1beta1.TFJob{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, bufSize).Decode(desiredJobSpec); err != nil {
		logger.Error(err, "Yaml decode error")
		return nil, err
	}

	// Set the default namespace.
	if desiredJobSpec.GetNamespace() == "" {
		desiredJobSpec.SetNamespace("default")
	}
	// Set the label for the job and pods.
	if desiredJobSpec.Labels == nil {
		desiredJobSpec.Labels = make(map[string]string)
	}
	desiredJobSpec.Labels[util.LabelTrial] = instance.Name
	for k := range desiredJobSpec.Spec.TFReplicaSpecs {
		if desiredJobSpec.Spec.TFReplicaSpecs[k].Template.Labels == nil {
			desiredJobSpec.Spec.TFReplicaSpecs[k].Template.Labels = make(map[string]string)
		}
		desiredJobSpec.Spec.TFReplicaSpecs[k].Template.Labels[util.LabelTrial] = instance.Name
	}

	return desiredJobSpec, nil
}

func (g *GeneralComposer) GetDesiredJobType(instance *trialv1alpha2.Trial) (string, error) {
	spec, err := g.getDesiredJobSpec(instance)
	if err != nil {
		return "", err
	}

	return spec.GetKind(), nil
}

func (g *GeneralComposer) getDesiredJobSpec(instance *trialv1alpha2.Trial) (*unstructured.Unstructured, error) {
	bufSize := 1024
	logger := log.WithValues("trial", types.NamespacedName{Name: instance.GetName(), Namespace: instance.GetNamespace()})
	buf := bytes.NewBufferString(instance.Spec.RunSpec)

	desiredJobSpec := &unstructured.Unstructured{}
	if err := k8syaml.NewYAMLOrJSONDecoder(buf, bufSize).Decode(desiredJobSpec); err != nil {
		logger.Error(err, "Yaml decode error")
		return nil, err
	}
	if desiredJobSpec.GetNamespace() == "" {
		desiredJobSpec.SetNamespace("default")
	}

	return desiredJobSpec, nil
}

func New() Composer {
	return &GeneralComposer{}
}
