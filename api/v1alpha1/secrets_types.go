/*
Copyright 2021.

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

// SecretsSpec defines the desired state of Secrets
type SecretsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Secrets. Edit secrets_types.go to remove/update
	SecretName string                 `json:"secretName,omitempty"`
	Cfn        []CloudformationOutput `json:"cfn"`
}

type CloudformationOutput struct {
	//+kubebuilder:validation:MinLength=1
	StackName string `json:"stackName"`

	//+kubebuilder:validation:MinLength=1
	KeyName string `json:"key"`

	//+kubebuilder:validation:MinLength=1
	OutputKey string `json:"outputKey"`

	Type string `json:"type,omitempty"`
}

// SecretsStatus defines the observed state of Secrets
type SecretsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Secrets is the Schema for the secrets API
type Secrets struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretsSpec   `json:"spec,omitempty"`
	Status SecretsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SecretsList contains a list of Secrets
type SecretsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Secrets `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Secrets{}, &SecretsList{})
}
