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

package v1

import (
	"context"
	_ "k8s.io/code-generator"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// Instance
// +k8s:openapi-gen=true
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec,omitempty"`
	Status InstanceStatus `json:"status,omitempty"`
}

// InstanceList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Instance `json:"items"`
}

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	Resource    Resource    `json:"resource,omitempty"`
	Environment Environment `json:"environment,omitempty"`
}

type Resource struct {
	CpuLimit  int `json:"cpuLimit,omitempty"`
	RamLimit  int `json:"ramLimit,omitempty"`
	DiskLimit int `json:"diskLimit,omitempty"`
}

type Environment struct {
	Owner string `json:"owner,omitempty"`
	Os    string `json:"os,omitempty"`
}

var _ resource.Object = &Instance{}
var _ resourcestrategy.Validater = &Instance{}

func (in *Instance) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *Instance) NamespaceScoped() bool {
	return false
}

func (in *Instance) New() runtime.Object {
	return &Instance{}
}

func (in *Instance) NewList() runtime.Object {
	return &InstanceList{}
}

func (in *Instance) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "boanlab.boanlab",
		Version:  "v1",
		Resource: "instances",
	}
}

func (in *Instance) IsStorageVersion() bool {
	return true
}

func (in *Instance) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &InstanceList{}

func (in *InstanceList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	InstanceID string     `json:"instanceID,omitempty"`
	Status     string     `json:"status,omitempty"`
	Snapshots  []Snapshot `json:"snapshots"`
}

type Snapshot struct {
	Name      string `json:"name,omitempty"`
	Generated string `json:"generated,omitempty"`
	Size      int    `json:"size,omitempty"`
}

func (in InstanceStatus) SubResourceName() string {
	return "status"
}

// Instance implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &Instance{}

func (in *Instance) GetStatus() resource.StatusSubResource {
	return in.Status
}

// InstanceStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &InstanceStatus{}

func (in InstanceStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*Instance).Status = in
}

var _ resource.ObjectWithArbitrarySubResource = &Instance{}

func (in *Instance) GetArbitrarySubResources() []resource.ArbitrarySubResource {

	return []resource.ArbitrarySubResource{
		// +kubebuilder:scaffold:subresource
		&InstanceSnapshot{},
	}
}
