/*
Copyright 2016 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Igusher attempts to authenticate a token to a known user.
type Igusher struct {
	metav1.TypeMeta
	// ObjectMeta fulfills the metav1.ObjectMetaAccessor interface so that the stock
	// REST handler paths work
	metav1.ObjectMeta

	// Spec holds information about the request being evaluated
	Spec IgusherSpec

	// Status is filled in by the server and indicates whether the request can be authenticated.
	Status IgusherStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// IgusherList is a collection of PodDisruptionBudgets.
type IgusherList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	Items []Igusher
}

// IgusherSpec is a description of the token authentication request.
type IgusherSpec struct {
	// Token is the opaque bearer token.
	Token string
}

// IgusherStatus is the result of the token authentication request.
// This type mirrors the authentication.Token interface
type IgusherStatus struct {
	// Authenticated indicates that the token was associated with a known user.
	Authenticated bool
}
