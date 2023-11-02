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

	validatorv1alpha1 "github.com/spectrocloud-labs/validator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AwsValidatorSpec defines the desired state of AwsValidator
type AwsValidatorSpec struct {
	Auth              AwsAuth            `json:"auth" yaml:"auth"`
	DefaultRegion     string             `json:"defaultRegion" yaml:"defaultRegion"`
	IamRoleRules      []IamRoleRule      `json:"iamRoleRules,omitempty" yaml:"iamRoleRules,omitempty"`
	IamUserRules      []IamUserRule      `json:"iamUserRules,omitempty" yaml:"iamUserRules,omitempty"`
	IamGroupRules     []IamGroupRule     `json:"iamGroupRules,omitempty" yaml:"iamGroupRules,omitempty"`
	IamPolicyRules    []IamPolicyRule    `json:"iamPolicyRules,omitempty" yaml:"iamPolicyRules,omitempty"`
	ServiceQuotaRules []ServiceQuotaRule `json:"serviceQuotaRules,omitempty" yaml:"serviceQuotaRules,omitempty"`
	TagRules          []TagRule          `json:"tagRules,omitempty" yaml:"tagRules,omitempty"`
}

type AwsAuth struct {
	// Option 1: lookup AWS creds from a secret
	SecretName string `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	// Option 2: specify a service account (EKS)
	ServiceAccountName string `json:"serviceAccountName,omitempty" yaml:"serviceAccountName,omitempty"`
}

type IamRoleRule struct {
	IamRoleName string           `json:"iamRoleName" yaml:"iamRoleName"`
	Policies    []PolicyDocument `json:"iamPolicies" yaml:"iamPolicies"`
}

func (r IamRoleRule) Name() string {
	return r.IamRoleName
}

func (r IamRoleRule) IAMPolicies() []PolicyDocument {
	return r.Policies
}

type IamUserRule struct {
	IamUserName string           `json:"iamUserName" yaml:"iamUserName"`
	Policies    []PolicyDocument `json:"iamPolicies" yaml:"iamPolicies"`
}

func (r IamUserRule) Name() string {
	return r.IamUserName
}

func (r IamUserRule) IAMPolicies() []PolicyDocument {
	return r.Policies
}

type IamGroupRule struct {
	IamGroupName string           `json:"iamGroupName" yaml:"iamGroupName"`
	Policies     []PolicyDocument `json:"iamPolicies" yaml:"iamPolicies"`
}

func (r IamGroupRule) Name() string {
	return r.IamGroupName
}

func (r IamGroupRule) IAMPolicies() []PolicyDocument {
	return r.Policies
}

type IamPolicyRule struct {
	IamPolicyARN string           `json:"iamPolicyArn" yaml:"iamPolicyArn"`
	Policies     []PolicyDocument `json:"iamPolicies" yaml:"iamPolicies"`
}

func (r IamPolicyRule) Name() string {
	return r.IamPolicyARN
}

func (r IamPolicyRule) IAMPolicies() []PolicyDocument {
	return r.Policies
}

type PolicyDocument struct {
	Name       string           `json:"name" yaml:"name"`
	Version    string           `json:"version" yaml:"version"`
	Statements []StatementEntry `json:"statements" yaml:"statements"`
}

type StatementEntry struct {
	Condition *Condition `json:"condition,omitempty" yaml:"condition,omitempty"`
	Effect    string     `json:"effect" yaml:"effect"`
	Actions   []string   `json:"actions" yaml:"actions"`
	Resources []string   `json:"resources" yaml:"resources"`
}

type Condition struct {
	Type   string   `json:"type" yaml:"type"`
	Key    string   `json:"key" yaml:"key"`
	Values []string `json:"values" yaml:"values"`
}

func (c *Condition) String() string {
	return fmt.Sprintf("%s: %s=%s", c.Type, c.Key, c.Values)
}

type ServiceQuotaRule struct {
	Region        string         `json:"region" yaml:"region"`
	ServiceCode   string         `json:"serviceCode" yaml:"serviceCode"`
	ServiceQuotas []ServiceQuota `json:"serviceQuotas" yaml:"serviceQuotas"`
}

type ServiceQuota struct {
	Name   string `json:"name" yaml:"name"`
	Buffer int    `json:"buffer" yaml:"buffer"`
}

type TagRule struct {
	Key           string   `json:"key" yaml:"key"`
	ExpectedValue string   `json:"expectedValue" yaml:"expectedValue"`
	Region        string   `json:"region" yaml:"region"`
	ResourceType  string   `json:"resourceType" yaml:"resourceType"`
	ARNs          []string `json:"arns" yaml:"arns"`
}

// AwsValidatorStatus defines the observed state of AwsValidator
type AwsValidatorStatus struct {
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []validatorv1alpha1.ValidationCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AwsValidator is the Schema for the awsvalidators API
type AwsValidator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AwsValidatorSpec   `json:"spec,omitempty"`
	Status AwsValidatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AwsValidatorList contains a list of AwsValidator
type AwsValidatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AwsValidator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AwsValidator{}, &AwsValidatorList{})
}
