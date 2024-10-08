//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Artifact) DeepCopyInto(out *Artifact) {
	*out = *in
	if in.ValidationType != nil {
		in, out := &in.ValidationType, &out.ValidationType
		*out = new(ValidationType)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Artifact.
func (in *Artifact) DeepCopy() *Artifact {
	if in == nil {
		return nil
	}
	out := new(Artifact)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Auth) DeepCopyInto(out *Auth) {
	*out = *in
	if in.SecretName != nil {
		in, out := &in.SecretName, &out.SecretName
		*out = new(string)
		**out = **in
	}
	if in.Basic != nil {
		in, out := &in.Basic, &out.Basic
		*out = new(BasicAuth)
		**out = **in
	}
	if in.ECR != nil {
		in, out := &in.ECR, &out.ECR
		*out = new(ECRAuth)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Auth.
func (in *Auth) DeepCopy() *Auth {
	if in == nil {
		return nil
	}
	out := new(Auth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BasicAuth) DeepCopyInto(out *BasicAuth) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BasicAuth.
func (in *BasicAuth) DeepCopy() *BasicAuth {
	if in == nil {
		return nil
	}
	out := new(BasicAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ECRAuth) DeepCopyInto(out *ECRAuth) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ECRAuth.
func (in *ECRAuth) DeepCopy() *ECRAuth {
	if in == nil {
		return nil
	}
	out := new(ECRAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OciRegistryRule) DeepCopyInto(out *OciRegistryRule) {
	*out = *in
	out.ManuallyNamed = in.ManuallyNamed
	if in.Artifacts != nil {
		in, out := &in.Artifacts, &out.Artifacts
		*out = make([]Artifact, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Auth.DeepCopyInto(&out.Auth)
	in.SignatureVerification.DeepCopyInto(&out.SignatureVerification)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OciRegistryRule.
func (in *OciRegistryRule) DeepCopy() *OciRegistryRule {
	if in == nil {
		return nil
	}
	out := new(OciRegistryRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OciValidator) DeepCopyInto(out *OciValidator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OciValidator.
func (in *OciValidator) DeepCopy() *OciValidator {
	if in == nil {
		return nil
	}
	out := new(OciValidator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OciValidator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OciValidatorList) DeepCopyInto(out *OciValidatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]OciValidator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OciValidatorList.
func (in *OciValidatorList) DeepCopy() *OciValidatorList {
	if in == nil {
		return nil
	}
	out := new(OciValidatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OciValidatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OciValidatorSpec) DeepCopyInto(out *OciValidatorSpec) {
	*out = *in
	if in.OciRegistryRules != nil {
		in, out := &in.OciRegistryRules, &out.OciRegistryRules
		*out = make([]OciRegistryRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OciValidatorSpec.
func (in *OciValidatorSpec) DeepCopy() *OciValidatorSpec {
	if in == nil {
		return nil
	}
	out := new(OciValidatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OciValidatorStatus) DeepCopyInto(out *OciValidatorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OciValidatorStatus.
func (in *OciValidatorStatus) DeepCopy() *OciValidatorStatus {
	if in == nil {
		return nil
	}
	out := new(OciValidatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SignatureVerification) DeepCopyInto(out *SignatureVerification) {
	*out = *in
	if in.PublicKeys != nil {
		in, out := &in.PublicKeys, &out.PublicKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SignatureVerification.
func (in *SignatureVerification) DeepCopy() *SignatureVerification {
	if in == nil {
		return nil
	}
	out := new(SignatureVerification)
	in.DeepCopyInto(out)
	return out
}
