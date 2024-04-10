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
// serverplatform.k8s.io API group.
package validation

import (
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/apis/serverplatform"
)

// ValidateServer.
func ValidateServer(s *serverplatform.Server) field.ErrorList {
	allErrs := field.ErrorList{}
	specPath := field.NewPath("spec")

	if !strings.HasPrefix(s.Spec.BuildLabel, "//") {
		allErrs = append(allErrs, field.Invalid(specPath.Child("BuildLabel"), s.Spec.BuildLabel, "BuildLabel must start with double-slash '//'"))
	}
	return allErrs
}

// ValidateBorgJob.
func ValidateBorgJob(bj *serverplatform.BorgJob) field.ErrorList {
	allErrs := field.ErrorList{}
	specPath := field.NewPath("spec")

	if bj.Spec.TaskCount <= 0 {
		allErrs = append(allErrs, field.Invalid(specPath.Child("TaskCount"), bj.Spec.TaskCount, "Size must be > 0"))
	}
	return allErrs
}

// ValidateEnvironment.
func ValidateEnvironment(env *serverplatform.Environment) field.ErrorList {
	allowedDataRealm := map[string]bool{
		"PROD": true,
		"TEST": true,
	}
	allErrs := field.ErrorList{}
	specPath := field.NewPath("spec")

	if _, ok := allowedDataRealm[env.Spec.DataRealm]; !ok {
		allErrs = append(allErrs, field.Invalid(specPath.Child("DataRealm"), env.Spec.DataRealm, "DataRealm must be one of: [PROD, Test]"))
	}
	return allErrs
}
