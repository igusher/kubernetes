/*
Copyright 2014 The Kubernetes Authors.
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

package server

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/serverplatform"
	"k8s.io/kubernetes/pkg/apis/serverplatform/validation"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// serverStrategy implements verification logic for Server.
type serverStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Server objects.
var Strategy = serverStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped returns true because all Servers need to be within a namespace.
func (serverStrategy) NamespaceScoped() bool {
	return true
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (serverStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"serverplatform/v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("status"),
		),
	}

	return fields
}

// PrepareForCreate clears the status of an Server before creation.
func (serverStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	server := obj.(*serverplatform.Server)
	// create cannot set status
	server.Status = serverplatform.ServerStatus{}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (serverStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newServer := obj.(*serverplatform.Server)
	oldServer := old.(*serverplatform.Server)
	// Update is not allowed to set status
	newServer.Status = oldServer.Status
}

// Validate validates a new Server.
func (serverStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	server := obj.(*serverplatform.Server)
	//return validation.ValidateServer(server, opts)
	return validation.ValidateServer(server)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (serverStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (serverStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is true for Server; this means you may create one with a PUT request.
func (serverStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (serverStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	//return validation.ValidateServer(obj.(*serverplatform.Server), opts)
	return validation.ValidateServer(obj.(*serverplatform.Server))
}

// WarningsOnUpdate returns warnings for the given update.
func (serverStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// AllowUnconditionalUpdate is the default update policy for Server objects. Status update should
// only be allowed if version match.
func (serverStrategy) AllowUnconditionalUpdate() bool {
	return false
}

type serverplatformStatusStrategy struct {
	serverStrategy
}

// StatusStrategy is the default logic invoked when updating object status.
var StatusStrategy = serverplatformStatusStrategy{Strategy}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (serverplatformStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"serverplatform/v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}

	return fields
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update of status
func (serverplatformStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newServer := obj.(*serverplatform.Server)
	oldServer := old.(*serverplatform.Server)
	// status changes are not allowed to update spec
	newServer.Spec = oldServer.Spec
}

// ValidateUpdate is the default update validation for an end user updating status
func (serverplatformStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (serverplatformStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
