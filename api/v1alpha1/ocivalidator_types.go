/*
Copyright 2023.

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

// OciValidatorSpec defines the desired state of OciValidator
type OciValidatorSpec struct {
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="OciRegistryRules must have a unique RuleName",rule="self.all(e, size(self.filter(x, x.name == e.name)) == 1)"
	OciRegistryRules []OciRegistryRule `json:"ociRegistryRules,omitempty" yaml:"ociRegistryRules,omitempty"`
}

// ResultCount returns the number of validation results expected for an OciValidatorSpec
func (s OciValidatorSpec) ResultCount() int {
	return len(s.OciRegistryRules)
}

// OciRegistryRule defines the validation rule for an OCI registry
type OciRegistryRule struct {
	// Name is the name of the rule
	RuleName string `json:"name" yaml:"name"`

	// Host is a reference to the host URL of an OCI compliant registry
	Host string `json:"host" yaml:"host"`

	// Artifacts is a slice of artifacts in the host registry that should be validated.
	Artifacts []Artifact `json:"artifacts,omitempty" yaml:"artifacts,omitempty"`

	// Auth provides authentication information for the registry
	Auth Auth `json:"auth,omitempty" yaml:"auth,omitempty"`

	// CaCert is the base64 encoded CA Certificate
	CaCert string `json:"caCert,omitempty" yaml:"caCert,omitempty"`

	// SignatureVerification provides the option to verify the signature of the image
	SignatureVerification SignatureVerification `json:"signatureVerification,omitempty" yaml:"signatureVerification,omitempty"`
}

// Name returns the name of the OciRegistryRule
func (r OciRegistryRule) Name() string {
	return r.RuleName
}

// Artifact defines the artifact to be validated
type Artifact struct {
	// Ref is the path to the artifact in the host registry that should be validated.
	// An individual artifact can take any of the following forms:
	// <repository-path>/<artifact-name>
	// <repository-path>/<artifact-name>:<tag>
	// <repository-path>/<artifact-name>@<digest>
	//
	// When no tag or digest are specified, the default tag "latest" is used.
	Ref string `json:"ref" yaml:"ref"`

	// LayerValidation specifies whether deep validation of the artifact layers should be performed.
	// The existence of layers is always validated, but this option allows for the deep validation of the layers.
	// See more details here:
	// https://github.com/google/go-containerregistry/blob/8dadbe76ff8c20d0e509406f04b7eade43baa6c1/pkg/v1/validate/image.go#L105
	LayerValidation bool `json:"layerValidation,omitempty" yaml:"layerValidation,omitempty"`
}

// Auth defines the authentication information for the registry
type Auth struct {
	// SecretName is the name of the Kubernetes Secret that exists in the same namespace as the OciValidator
	// and that contains the credentials used to authenticate to the OCI Registry
	SecretName string `json:"secretName" yaml:"secretName"`
}

// SignatureVerification defines the provider and secret name to verify the signatures of artifacts in an OCI registry
type SignatureVerification struct {
	// Provider specifies the technology used to sign the OCI Artifact
	// +kubebuilder:validation:Enum=cosign
	// +kubebuilder:default:=cosign
	Provider string `json:"provider" yaml:"provider"`

	// SecretName is the name of the Kubernetes Secret that exists in the same namespace as the OciValidator
	// and that contains the trusted public keys used to sign artifacts in the OciRegistryRule
	SecretName string `json:"secretName" yaml:"secretName"`
}

// OciValidatorStatus defines the observed state of OciValidator
type OciValidatorStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OciValidator is the Schema for the ocivalidators API
type OciValidator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OciValidatorSpec   `json:"spec,omitempty"`
	Status OciValidatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OciValidatorList contains a list of OciValidator
type OciValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OciValidator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OciValidator{}, &OciValidatorList{})
}
