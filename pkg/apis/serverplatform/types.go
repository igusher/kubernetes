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

package serverplatform

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Server struct {
	metav1.TypeMeta
	// ObjectMeta fulfills the metav1.ObjectMetaAccessor interface so that the stock
	// REST handler paths work
	metav1.ObjectMeta
	Spec ServerSpec
	Status ServerStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ServerList
type ServerList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	Items []Server
}

// ServerSPec
type ServerSpec struct {
  Name string
  BuildLabel string
  GanpatiPrefix string
  GSLBPrefix string
  OwnerTeam string
}

// ServerStatus
type ServerStatus struct {
	Reconciling bool
	Etag string
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BorgJob struct {
	metav1.TypeMeta
	// ObjectMeta fulfills the metav1.ObjectMetaAccessor interface so that the stock
	// REST handler paths work
	metav1.ObjectMeta
	Spec BorgJobSpec
	Status BorgJobStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// BorgJobList
type BorgJobList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	Items []BorgJob
}

// BorgJobSpec
type BorgJobSpec struct {
  Environment string
  Cell string
  TaskCount int32
}

// BorgJobStatus
type BorgJobStatus struct {
	Reconciling bool
	Etag string
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Environment struct {
	metav1.TypeMeta
	// ObjectMeta fulfills the metav1.ObjectMetaAccessor interface so that the stock
	// REST handler paths work
	metav1.ObjectMeta
	Spec EnvironmentSpec
	Status EnvironmentStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// EnvironmentList
type EnvironmentList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta
	Items []Environment
}

// EnvironmentSpec
type EnvironmentSpec struct {
  DataRealm 	     string
  ReleaseStage 	     string
  RuntimeEnvironment string
}

// EnvironemntStatus
type EnvironmentStatus struct {
	Reconciling bool
	Etag string
}

