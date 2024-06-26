//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ACommand) DeepCopyInto(out *ACommand) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ACommand.
func (in *ACommand) DeepCopy() *ACommand {
	if in == nil {
		return nil
	}
	out := new(ACommand)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ACommandRan) DeepCopyInto(out *ACommandRan) {
	*out = *in
	if in.CommandStatus != nil {
		in, out := &in.CommandStatus, &out.CommandStatus
		*out = make([]ACommandStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ACommandRan.
func (in *ACommandRan) DeepCopy() *ACommandRan {
	if in == nil {
		return nil
	}
	out := new(ACommandRan)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ACommandStatus) DeepCopyInto(out *ACommandStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ACommandStatus.
func (in *ACommandStatus) DeepCopy() *ACommandStatus {
	if in == nil {
		return nil
	}
	out := new(ACommandStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AContainerNames) DeepCopyInto(out *AContainerNames) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AContainerNames.
func (in *AContainerNames) DeepCopy() *AContainerNames {
	if in == nil {
		return nil
	}
	out := new(AContainerNames)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AEnv) DeepCopyInto(out *AEnv) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AEnv.
func (in *AEnv) DeepCopy() *AEnv {
	if in == nil {
		return nil
	}
	out := new(AEnv)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AMappingConfig) DeepCopyInto(out *AMappingConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AMappingConfig.
func (in *AMappingConfig) DeepCopy() *AMappingConfig {
	if in == nil {
		return nil
	}
	out := new(AMappingConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AMountOptions) DeepCopyInto(out *AMountOptions) {
	*out = *in
	if in.MappingConfig != nil {
		in, out := &in.MappingConfig, &out.MappingConfig
		*out = make([]AMappingConfig, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AMountOptions.
func (in *AMountOptions) DeepCopy() *AMountOptions {
	if in == nil {
		return nil
	}
	out := new(AMountOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APodExecutionContext) DeepCopyInto(out *APodExecutionContext) {
	*out = *in
	in.CurrentStep.DeepCopyInto(&out.CurrentStep)
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]AStep, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APodExecutionContext.
func (in *APodExecutionContext) DeepCopy() *APodExecutionContext {
	if in == nil {
		return nil
	}
	out := new(APodExecutionContext)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AStep) DeepCopyInto(out *AStep) {
	*out = *in
	if in.Commands != nil {
		in, out := &in.Commands, &out.Commands
		*out = make([]ACommand, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Environments.DeepCopyInto(&out.Environments)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AStep.
func (in *AStep) DeepCopy() *AStep {
	if in == nil {
		return nil
	}
	out := new(AStep)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AStepEnv) DeepCopyInto(out *AStepEnv) {
	*out = *in
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]AEnv, len(*in))
		copy(*out, *in)
	}
	in.MountOptions.DeepCopyInto(&out.MountOptions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AStepEnv.
func (in *AStepEnv) DeepCopy() *AStepEnv {
	if in == nil {
		return nil
	}
	out := new(AStepEnv)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AStepPodInfo) DeepCopyInto(out *AStepPodInfo) {
	*out = *in
	if in.ConatinerNames != nil {
		in, out := &in.ConatinerNames, &out.ConatinerNames
		*out = make([]AContainerNames, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AStepPodInfo.
func (in *AStepPodInfo) DeepCopy() *AStepPodInfo {
	if in == nil {
		return nil
	}
	out := new(AStepPodInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AbsurdCI) DeepCopyInto(out *AbsurdCI) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AbsurdCI.
func (in *AbsurdCI) DeepCopy() *AbsurdCI {
	if in == nil {
		return nil
	}
	out := new(AbsurdCI)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AbsurdCI) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AbsurdCIList) DeepCopyInto(out *AbsurdCIList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AbsurdCI, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AbsurdCIList.
func (in *AbsurdCIList) DeepCopy() *AbsurdCIList {
	if in == nil {
		return nil
	}
	out := new(AbsurdCIList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AbsurdCIList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AbsurdCISpec) DeepCopyInto(out *AbsurdCISpec) {
	*out = *in
	if in.Steps != nil {
		in, out := &in.Steps, &out.Steps
		*out = make([]AStep, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AbsurdCISpec.
func (in *AbsurdCISpec) DeepCopy() *AbsurdCISpec {
	if in == nil {
		return nil
	}
	out := new(AbsurdCISpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AbsurdCIStatus) DeepCopyInto(out *AbsurdCIStatus) {
	*out = *in
	in.APodExecutionContextInfo.DeepCopyInto(&out.APodExecutionContextInfo)
	if in.AStepPodCreationInfo != nil {
		in, out := &in.AStepPodCreationInfo, &out.AStepPodCreationInfo
		*out = make(map[string]AStepPodInfo, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Dag != nil {
		in, out := &in.Dag, &out.Dag
		*out = make([]AStep, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AbsurdCIStatus.
func (in *AbsurdCIStatus) DeepCopy() *AbsurdCIStatus {
	if in == nil {
		return nil
	}
	out := new(AbsurdCIStatus)
	in.DeepCopyInto(out)
	return out
}
