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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OciValidatorSpec defines the desired state of OciValidator
type OciValidatorSpec struct {
	OciRegistryRules []OciRegistryRule `json:"ociRegistryRules,omitempty" yaml:"ociRegistryRules,omitempty"`
	EcrRegistryRules []EcrRegistryRule `json:"ecrRegistryRules,omitempty" yaml:"ecrRegistryRules,omitempty"`
}

func (s OciValidatorSpec) ResultCount() int {
	return len(s.EcrRegistryRules) + len(s.OciRegistryRules)
}

type OciRegistryRule struct {
	Host            string    `json:"host" yaml:"host"`
	RepositoryPaths []string  `json:"repositoryPaths,omitempty" yaml:"repositoryPaths,omitempty"`
	Auth            BasicAuth `json:"auth,omitempty" yaml:"auth,omitempty"`
	Cert            string    `json:"cert,omitempty" yaml:"cert,omitempty"`
}

func (r OciRegistryRule) Name() string {
	return fmt.Sprintf("%s/%s", r.Host, r.RepositoryPaths)
}

type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type EcrRegistryRule struct {
	Host            string   `json:"host" yaml:"host"`
	RepositoryPaths []string `json:"repositoryPaths,omitempty" yaml:"repositoryPaths,omitempty"`
	Auth            EcrAuth  `json:"auth,omitempty" yaml:"auth,omitempty"`
}

func (r EcrRegistryRule) Name() string {
	return fmt.Sprintf("%s/%s", r.Host, r.RepositoryPaths)
}

type EcrAuth struct {
	RoleArn        string         `json:"roleArn,omitempty" yaml:"roleArn,omitempty"`
	EcrCredentials EcrCredentials `json:"ecrCredentials,omitempty" yaml:"ecrCredentials,omitempty"`
}

type EcrCredentials struct {
	AccessKey       string `json:"accessKey" yaml:"accessKey"`
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey"`
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
