/*
Copyright 2018 The Kubernetes Authors.

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

// Package validation contains methods to validate kinds in the
// igusher.k8s.io API group.
package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/apis/igusher"
)

// ValidateIgusher validates a TokenRequest.
func ValidateIgusher(tr *igusher.Igusher) field.ErrorList {
	allErrs := field.ErrorList{}
	specPath := field.NewPath("spec")

	if tr.Spec.Token == "trigger" {
		allErrs = append(allErrs, field.Invalid(specPath.Child("Token"), tr.Spec.Token, "may not specify Token='trigger'"))
	}
	return allErrs
}
