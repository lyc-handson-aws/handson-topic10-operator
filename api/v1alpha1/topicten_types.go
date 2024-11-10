/*
Copyright 2024.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TopicTenSpec defines the desired state of TopicTen
type TopicTenSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=5
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +kubebuilder:validation:Required
	TargetArn string `json:"targetarn"`

	// +kubebuilder:validation:Required
	KMSArn string `json:"kmsarn"`

	// +kubebuilder:validation:Required
	CloudWatchArn string `json:"cloudwatcharn"`
}

// TopicTenStatus defines the observed state of TopicTen
type TopicTenStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CurrentReplicas int32 `json:"currentReplicas"`

	//Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TopicTen is the Schema for the topictens API
type TopicTen struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopicTenSpec   `json:"spec,omitempty"`
	Status TopicTenStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TopicTenList contains a list of TopicTen
type TopicTenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TopicTen `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TopicTen{}, &TopicTenList{})
}
