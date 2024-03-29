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
	"fmt"

	_ "k8s.io/code-generator"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcerest"
	contextutil "sigs.k8s.io/apiserver-runtime/pkg/util/context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

var _ resource.GetterUpdaterSubResource = &InstanceSnapshot{}
var _ resourcerest.Getter = &InstanceSnapshot{}
var _ resourcerest.Updater = &InstanceSnapshot{}

// InstanceSnapshot
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type InstanceSnapshot struct {
	metav1.TypeMeta `json:",inline" `
}

func (c *InstanceSnapshot) Destroy() {}

func (c *InstanceSnapshot) SubResourceName() string {
	return "snapshot"
}

func (c *InstanceSnapshot) New() runtime.Object {
	return &InstanceSnapshot{}
}

func (c *InstanceSnapshot) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	// EDIT IT
	parentStorage, ok := contextutil.GetParentStorage(ctx)
	if !ok {
		return nil, fmt.Errorf("no parent storage found in the context")
	}
	return parentStorage.Get(ctx, name, options)
}

func (c *InstanceSnapshot) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// EDIT IT
	parentStorage, ok := contextutil.GetParentStorage(ctx)
	if !ok {
		return nil, false, fmt.Errorf("no parent storage found in the context")
	}
	return parentStorage.Update(ctx, name, objInfo, createValidation, updateValidation, forceAllowCreate, options)
}
