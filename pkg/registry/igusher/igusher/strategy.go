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

package igusher

import (
	"context"

	//apiequality "k8s.io/apimachinery/pkg/api/equality"
	//metav1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	//genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/storage/names"
	//utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/igusher"
	"k8s.io/kubernetes/pkg/apis/igusher/validation"
	//"k8s.io/kubernetes/pkg/features"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// igusherStrategy implements verification logic for PodDisruptionBudgets.
type igusherStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating PodDisruptionBudget objects.
var Strategy = igusherStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped returns true because all PodDisruptionBudget' need to be within a namespace.
func (igusherStrategy) NamespaceScoped() bool {
	return true
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (igusherStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"igusher/v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("status"),
		),
	}

	return fields
}

// PrepareForCreate clears the status of an PodDisruptionBudget before creation.
func (igusherStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	podDisruptionBudget := obj.(*igusher.Igusher)
	// create cannot set status
	podDisruptionBudget.Status = igusher.IgusherStatus{}

	//podDisruptionBudget.Generation = 1

	//dropDisabledFields(&podDisruptionBudget.Spec, nil)
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (igusherStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newPodDisruptionBudget := obj.(*igusher.Igusher)
	oldPodDisruptionBudget := old.(*igusher.Igusher)
	// Update is not allowed to set status
	newPodDisruptionBudget.Status = oldPodDisruptionBudget.Status

	// Any changes to the spec increment the generation number, any changes to the
	// status should reflect the generation number of the corresponding object.
	// See metav1.ObjectMeta description for more information on Generation.
	//if !apiequality.Semantic.DeepEqual(oldPodDisruptionBudget.Spec, newPodDisruptionBudget.Spec) {
	//	newPodDisruptionBudget.Generation = oldPodDisruptionBudget.Generation + 1
	//}

	//dropDisabledFields(&newPodDisruptionBudget.Spec, &oldPodDisruptionBudget.Spec)
}

// Validate validates a new PodDisruptionBudget.
func (igusherStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	podDisruptionBudget := obj.(*igusher.Igusher)
	//return validation.ValidateIgusher(podDisruptionBudget, opts)
	return validation.ValidateIgusher(podDisruptionBudget)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (igusherStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
func (igusherStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is true for PodDisruptionBudget; this means you may create one with a PUT request.
func (igusherStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (igusherStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	//return validation.ValidateIgusher(obj.(*igusher.Igusher), opts)
	return validation.ValidateIgusher(obj.(*igusher.Igusher))
}

// WarningsOnUpdate returns warnings for the given update.
func (igusherStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

// AllowUnconditionalUpdate is the default update policy for PodDisruptionBudget objects. Status update should
// only be allowed if version match.
func (igusherStrategy) AllowUnconditionalUpdate() bool {
	return false
}

type igusherStatusStrategy struct {
	igusherStrategy
}

// StatusStrategy is the default logic invoked when updating object status.
var StatusStrategy = igusherStatusStrategy{Strategy}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
func (igusherStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"igusher/v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}

	return fields
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update of status
func (igusherStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newPodDisruptionBudget := obj.(*igusher.Igusher)
	oldPodDisruptionBudget := old.(*igusher.Igusher)
	// status changes are not allowed to update spec
	newPodDisruptionBudget.Spec = oldPodDisruptionBudget.Spec
}

// ValidateUpdate is the default update validation for an end user updating status
func (igusherStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	/*var apiVersion schema.GroupVersion
	if requestInfo, found := genericapirequest.RequestInfoFrom(ctx); found {
		apiVersion = schema.GroupVersion{
			Group:   requestInfo.APIGroup,
			Version: requestInfo.APIVersion,
		}
	}
	return validation.ValidatePodDisruptionBudgetStatusUpdate(obj.(*igusher.Igusher).Status,
		old.(*igusher.Igusher).Status, field.NewPath("status"), apiVersion)*/
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (igusherStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

/*
// dropDisabledFields removes disabled fields from the pod disruption budget spec.
// This should be called from PrepareForCreate/PrepareForUpdate for all resources containing a pod disruption budget spec.
func dropDisabledFields(pdbSpec, oldPDBSpec *igusher.IgusherSpec) {
	if !utilfeature.DefaultFeatureGate.Enabled(features.PDBUnhealthyPodEvictionPolicy) {
		if !unhealthyPodEvictionPolicyInUse(oldPDBSpec) {
			pdbSpec.UnhealthyPodEvictionPolicy = nil
		}
	}
}

func unhealthyPodEvictionPolicyInUse(oldPDBSpec *igusher.IgusherSpec) bool {
	if oldPDBSpec == nil {
		return false
	}
	if oldPDBSpec.UnhealthyPodEvictionPolicy != nil {
		return true
	}
	return false
}*/
