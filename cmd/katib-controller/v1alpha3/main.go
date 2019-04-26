/*
Copyright 2018 The Kubeflow Authors

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

/*
 Katib-controller is a controller (operator) for Experiments and Trials
*/
package main

import (
	"log"

	"github.com/kubeflow/katib/pkg/api/operators/apis"
	controller "github.com/kubeflow/katib/pkg/controller/v1alpha3"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	logf.SetLogger(logf.ZapLogger(false))
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("config.GetConfig()")
		log.Fatal(err)
	}

	// Create a new katib controller to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		log.Printf("manager.New")
		log.Fatal(err)
	}

	log.Printf("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Printf("apis.AddToScheme")
		log.Fatal(err)
	}

	// Setup katib controller
	if err := controller.AddToManager(mgr); err != nil {
		log.Printf("controller.AddToManager(mgr)")
		log.Fatal(err)
	}

	log.Printf("Starting the Cmd.")

	// Starting the katib controller
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
