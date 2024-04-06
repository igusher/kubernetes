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

package borgjob

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

// borgjobStrategy implements verification logic for BorgJob.
type borgjobStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating BorgJob objects.
var Strategy = borgjobStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped returns true because all BorgJobs need to be within a namespace.
func (borgjobStrategy) NamespaceScoped() bool {
	return true
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (borgjobStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"serverplatform/v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("status"),
		),
	}

	return fields
}

// PrepareForCreate clears the status of an BorgJob before creation.
func (borgjobStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	server := obj.(*serverplatform.BorgJob)
	// create cannot set status
	server.Status = serverplatform.BorgJobStatus{}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (borgjobStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newBorgJob := obj.(*serverplatform.BorgJob)
	oldBorgJob := old.(*serverplatform.BorgJob)
	// Update is not allowed to set status
	newBorgJob.Status = oldBorgJob.Status
}

// Validate validates a new BorgJob.
func (borgjobStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	server := obj.(*serverplatform.BorgJob)
	//return validation.ValidateBorgJob(server, opts)
	return validation.ValidateBorgJob(server)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (borgjobStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (borgjobStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is true for BorgJob; this means you may create one with a PUT request.
func (borgjobStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (borgjobStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	//return validation.ValidateBorgJob(obj.(*serverplatform.BorgJob), opts)
	return validation.ValidateBorgJob(obj.(*serverplatform.BorgJob))
}

// WarningsOnUpdate returns warnings for the given update.
func (borgjobStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// AllowUnconditionalUpdate is the default update policy for BorgJob objects. Status update should
// only be allowed if version match.
func (borgjobStrategy) AllowUnconditionalUpdate() bool {
	return false
}

type serverplatformStatusStrategy struct {
	borgjobStrategy
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
	newBorgJob := obj.(*serverplatform.BorgJob)
	oldBorgJob := old.(*serverplatform.BorgJob)
	// status changes are not allowed to update spec
	newBorgJob.Spec = oldBorgJob.Spec
}

// ValidateUpdate is the default update validation for an end user updating status
func (serverplatformStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (serverplatformStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
