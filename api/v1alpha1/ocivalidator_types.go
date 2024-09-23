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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/validator-labs/validator/pkg/plugins"
	"github.com/validator-labs/validator/pkg/validationrule"

	"github.com/validator-labs/validator-plugin-oci/pkg/constants"
)

// OciValidatorSpec defines the desired state of OciValidator.
type OciValidatorSpec struct {
	// +kubebuilder:validation:MaxItems=5
	// +kubebuilder:validation:XValidation:message="OciRegistryRules must have a unique RuleName",rule="self.all(e, size(self.filter(x, x.name == e.name)) == 1)"
	OciRegistryRules []OciRegistryRule `json:"ociRegistryRules,omitempty" yaml:"ociRegistryRules,omitempty"`
}

var _ plugins.PluginSpec = (*OciValidatorSpec)(nil)

// PluginCode returns the OCI validator's plugin code.
func (s OciValidatorSpec) PluginCode() string {
	return constants.PluginCode
}

// ResultCount returns the number of validation results expected for an OciValidatorSpec.
func (s OciValidatorSpec) ResultCount() int {
	return len(s.OciRegistryRules)
}

// OciRegistryRule defines the validation rule for an OCI registry.
type OciRegistryRule struct {
	validationrule.ManuallyNamed `json:",inline" yaml:",omitempty"`

	// Name is a unique name for the OciRegistryRule.
	RuleName string `json:"name" yaml:"name"`

	// Host is the URI of an OCI registry.
	Host string `json:"host" yaml:"host"`

	// ValidationType specifies which (if any) type of validation is performed on the artifacts.
	// Valid values are "full", "fast", and "none". When set to "none", the artifact will not be pulled and no extra validation will be performed.
	// For both "full" and "fast" validationType, the following validations will be executed:
	// - Layers existence will be validated
	// - Config digest, size, content, and type will be validated
	// - Manifest digest, content, and size will be validated
	// For "full" validationType, the following additional validations will be performed:
	// - Layer digest, diffID, size, and media type will be validated
	// See more details about validation here:
	// https://github.com/google/go-containerregistry/blob/8dadbe76ff8c20d0e509406f04b7eade43baa6c1/pkg/v1/validate/image.go#L30
	// +kubebuilder:validation:Enum=full;fast;none
	// +kubebuilder:default:=none
	ValidationType ValidationType `json:"validationType" yaml:"validationType"`

	// Artifacts is a slice of artifacts in the OCI registry that should be validated.
	// +kubebuilder:validation:MinItems=1
	Artifacts []Artifact `json:"artifacts,omitempty" yaml:"artifacts,omitempty"`

	// Auth provides authentication information for the registry.
	Auth Auth `json:"auth,omitempty" yaml:"auth,omitempty"`

	// InsecureSkipTLSVerify specifies whether to skip verification of the OCI registry's TLS certificate.
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty" yaml:"insecureSkipTLSVerify,omitempty"`

	// CaCert is the CA certificate of the OCI registry.
	CaCert string `json:"caCert,omitempty" yaml:"caCert,omitempty"`

	// SignatureVerification provides signature verification options for the artifacts.
	SignatureVerification SignatureVerification `json:"signatureVerification,omitempty" yaml:"signatureVerification,omitempty"`
}

var _ validationrule.Interface = (*OciRegistryRule)(nil)

// ValidationType defines the type of extra validation to perform on the artifacts.
type ValidationType string

const (
	// ValidationTypeFull specifies full validation of the artifacts.
	ValidationTypeFull ValidationType = "full"
	// ValidationTypeFast specifies fast validation of the artifacts.
	ValidationTypeFast ValidationType = "fast"
	// ValidationTypeNone specifies no extra validation of the artifacts, artifacts will not be pulled.
	ValidationTypeNone ValidationType = "none"
)

// Name returns the name of the OciRegistryRule.
func (r OciRegistryRule) Name() string {
	return r.RuleName
}

// SetName sets the name of the OciRegistryRule.
func (r *OciRegistryRule) SetName(name string) {
	r.RuleName = name
}

// Artifact defines an OCI artifact to be validated.
type Artifact struct {
	// Ref is the path to the artifact in the host registry that should be validated.
	// An individual artifact can take any of the following forms:
	// <repository-path>/<artifact-name>
	// <repository-path>/<artifact-name>:<tag>
	// <repository-path>/<artifact-name>@<digest>
	//
	// When no tag or digest are specified, the default tag "latest" is used.
	Ref string `json:"ref" yaml:"ref"`

	// ValidationType overrides the OciRegistryRule level ValidationType for a particular artifact.
	// +kubebuilder:validation:Enum=full;fast;none
	ValidationType *ValidationType `json:"validationType,omitempty" yaml:"validationType,omitempty"`
}

// Auth defines the authentication information for the registry.
// One of SecretName, Basic, or ECR must be provided for a private registry.
// If multiple fields are provided, the order of precedence is SecretName, Basic, ECR.
type Auth struct {
	// SecretName is the name of the Kubernetes Secret that exists in the same namespace as the OciValidator
	// and that contains the credentials used to authenticate to the OCI Registry.
	SecretName *string `json:"secretName,omitempty" yaml:"secretName,omitempty"`

	// BasicAuth is the username and password used to authenticate to the OCI registry.
	Basic *BasicAuth `json:"basic,omitempty" yaml:"basic,omitempty"`

	// ECRAuth is the access key ID, secret access key, and session token used to authenticate to ECR.
	ECR *ECRAuth `json:"ecr,omitempty" yaml:"ecr,omitempty"`
}

// BasicAuth defines the username and password used to authenticate to the OCI registry.
type BasicAuth struct {
	// Username is the username used to authenticate to the OCI Registry.
	Username string `json:"username" yaml:"username"`

	// Password is the password used to authenticate to the OCI Registry.
	Password string `json:"password" yaml:"password"`
}

// ECRAuth defines the access key ID, secret access key, and session token used to authenticate to ECR.
type ECRAuth struct {
	// AccessKeyID is the AWS access key ID used to authenticate to ECR.
	AccessKeyID string `json:"accessKeyID" yaml:"accessKeyID"`

	// SecretAccessKey is the AWS secret access key used to authenticate to ECR.
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey"`

	// SessionToken is the AWS session token used to authenticate to ECR.
	SessionToken string `json:"sessionToken,omitempty" yaml:"sessionToken,omitempty"`
}

// SignatureVerification defines the provider and secret name to verify the signatures of artifacts in an OCI registry.
type SignatureVerification struct {
	// Provider specifies the technology used to sign the OCI Artifact.
	// +kubebuilder:validation:Enum=cosign
	// +kubebuilder:default:=cosign
	Provider string `json:"provider" yaml:"provider"`

	// SecretName is the name of the Kubernetes Secret that exists in the same namespace as the OciValidator
	// and that contains the trusted public keys used to sign artifacts in the OciRegistryRule.
	SecretName string `json:"secretName" yaml:"secretName"`

	// PublicKeys is a slice of public keys used to verify the signatures of artifacts in the OciRegistryRule.
	PublicKeys []string `json:"publicKeys,omitempty" yaml:"publicKeys,omitempty"`
}

// OciValidatorStatus defines the observed state of OciValidator.
type OciValidatorStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OciValidator is the Schema for the ocivalidators API.
type OciValidator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OciValidatorSpec   `json:"spec,omitempty"`
	Status OciValidatorStatus `json:"status,omitempty"`
}

// GetKind returns the OCI validator's kind.
func (v OciValidator) GetKind() string {
	return reflect.TypeOf(v).Name()
}

// PluginCode returns the OCI validator's plugin code.
func (v OciValidator) PluginCode() string {
	return v.Spec.PluginCode()
}

// ResultCount returns the number of validation results expected for an OciValidator.
func (v OciValidator) ResultCount() int {
	return v.Spec.ResultCount()
}

//+kubebuilder:object:root=true

// OciValidatorList contains a list of OciValidator.
type OciValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OciValidator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OciValidator{}, &OciValidatorList{})
}

// BasicAuthsDirect returns a map of basic authentication details for each rule when invoked directly.
func (s *OciValidatorSpec) BasicAuthsDirect() map[string][]string {
	auths := make(map[string][]string)

	for _, r := range s.OciRegistryRules {
		if r.Auth.Basic != nil {
			auths[r.Name()] = []string{r.Auth.Basic.Username, r.Auth.Basic.Password}
			continue
		}
	}

	return auths
}

// AllPubKeysDirect returns a map of public keys for each rule when invoked directly.
func (s *OciValidatorSpec) AllPubKeysDirect() map[string][][]byte {
	pubKeysMap := make(map[string][][]byte)

	for _, r := range s.OciRegistryRules {
		if len(r.SignatureVerification.PublicKeys) == 0 {
			continue
		}

		pubKeys := make([][]byte, len(r.SignatureVerification.PublicKeys))
		for i, pk := range r.SignatureVerification.PublicKeys {
			pubKeys[i] = []byte(pk)
		}
		pubKeysMap[r.Name()] = pubKeys
	}

	return pubKeysMap
}
