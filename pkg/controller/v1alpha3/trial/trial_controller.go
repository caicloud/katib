/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trial

import (
	"context"

	tfv1beta1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/recorder"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/trial/clientset"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/trial/composer"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/trial/metricscollector"
)

const (
	controllerName = "trial-controller"
)

var log = logf.Log.WithName(controllerName)

// Add creates a new Trial Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	r := &ReconcileTrial{
		Client:   mgr.GetClient(),
		Recorder: recorder.New(mgr.GetRecorder(controllerName)),
		Composer: composer.New(),
		scheme:   mgr.GetScheme(),
	}

	mc, err := metricscollector.New()
	if err != nil {
		panic(err)
	}

	r.MetricsCollector = mc
	r.Unstructured = clientset.NewUnstructured(r.Client, r.Recorder)
	r.TensorFlow = clientset.NewTF(r.Client, r.Recorder)
	r.updateStatusHandler = r.updateStatus
	return r
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Trial
	err = c.Watch(&source.Kind{Type: &trialv1alpha2.Trial{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &tfv1beta1.TFJob{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &trialv1alpha2.Trial{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileTrial{}

// ReconcileTrial reconciles a Trial object
type ReconcileTrial struct {
	client.Client
	recorder.Recorder
	composer.Composer
	clientset.TensorFlow
	metricscollector.MetricsCollector
	clientset.Unstructured
	scheme *runtime.Scheme
	// To allow injection of updateStatus for testing.
	updateStatusHandler func(r *trialv1alpha2.Trial) error
}

// Reconcile reads that state of the cluster for a Trial object and makes changes based on the state read
// and what is in the Trial.Spec
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=trials.kubeflow.org,resources=trials,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=trials.kubeflow.org,resources=trials/status,verbs=get;update;patch
func (r *ReconcileTrial) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Trial instance
	instance := &trialv1alpha2.Trial{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return r.handle(instance)
}
