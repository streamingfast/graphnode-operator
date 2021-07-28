// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestRun) DeepCopyInto(out *TestRun) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestRun.
func (in *TestRun) DeepCopy() *TestRun {
	if in == nil {
		return nil
	}
	out := new(TestRun)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TestRun) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestRunList) DeepCopyInto(out *TestRunList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TestRun, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestRunList.
func (in *TestRunList) DeepCopy() *TestRunList {
	if in == nil {
		return nil
	}
	out := new(TestRunList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TestRunList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestRunSpec) DeepCopyInto(out *TestRunSpec) {
	*out = *in
	out.PGDataStorage = in.PGDataStorage.DeepCopy()
	out.OutputDirSize = in.OutputDirSize.DeepCopy()
	in.PostgresResources.DeepCopyInto(&out.PostgresResources)
	in.GraphnodeResources.DeepCopyInto(&out.GraphnodeResources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestRunSpec.
func (in *TestRunSpec) DeepCopy() *TestRunSpec {
	if in == nil {
		return nil
	}
	out := new(TestRunSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TestRunStatus) DeepCopyInto(out *TestRunStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TestRunStatus.
func (in *TestRunStatus) DeepCopy() *TestRunStatus {
	if in == nil {
		return nil
	}
	out := new(TestRunStatus)
	in.DeepCopyInto(out)
	return out
}
