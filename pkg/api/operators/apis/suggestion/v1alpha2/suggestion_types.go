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

package v1alpha2

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SuggestionSpec defines the desired state of Suggestion
type SuggestionSpec struct {
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Type of the suggestion.
	Type SuggestionType `json:"type,omitempty"`

	// Defines if the suggestion needs previous results.
	NeedHistory bool `json:"needHistory,omitempty"`

	// Template describes the pods that will be created.
	Template v1.PodTemplateSpec `json:"template"`
}

type SuggestionType string

const (
	SuggestionTypeNAS           = "NAS"
	SuggestionTypeEarlyStopping = "EarlyStopping"
	SuggestionTypeHPTuning      = "HyperParameter"
)

// SuggestionStatus defines the observed state of Suggestion
type SuggestionStatus struct {
	// Represents time when the Experiment was acknowledged by the Experiment controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the Experiment was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the Experiment was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	// List of observed runtime conditions for this Experiment.
	Conditions []SuggestionCondition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen=true
// SuggestionCondition describes the state of the suggestion at a certain point.
type SuggestionCondition struct {
	// Type of experiment condition.
	Type SuggestionConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`

	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`

	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`

	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// SuggestionConditionType defines the state of an Suggestion.
type SuggestionConditionType string

const (
	SuggestionDeploymentAvailable      SuggestionConditionType = "DeploymentAvailable"
	SuggestionDeploymentProgressing    SuggestionConditionType = "DeploymentProgressing"
	SuggestionDeploymentReplicaFailure SuggestionConditionType = "DeploymentReplicaFailure"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Suggestion is the Schema for the suggestions API
// +k8s:openapi-gen=true
type Suggestion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SuggestionSpec   `json:"spec,omitempty"`
	Status SuggestionStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SuggestionList contains a list of Suggestion
type SuggestionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Suggestion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Suggestion{}, &SuggestionList{})
}
