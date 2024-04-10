/*
Copyright 2017 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Server is my custom resource
type Server struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec ServerSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status ServerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Server
	Items []Server `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ServerSpec is a description of the token authentication request.
type ServerSpec struct {
	// +optional
	Name          string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	BuildLabel    string `json:"build_label,omitempty" protobuf:"bytes,2,opt,name=build_label"`
	GanpatiPrefix string `json:"ganpati_prefix,omitempty" protobuf:"bytes,3,opt,name=ganpati_prefix"`
	GSLBPrefix    string `json:"gslb_prefix,omitempty" protobuf:"bytes,4,opt,name=gslb_prefix"`
	OwnerTeam     string `json:"owner_team,omitempty" protobuf:"bytes,5,opt,name=owner_team"`
}

// ServerStatus is the result of the token authentication request.
type ServerStatus struct {
	// +optional
	Reconciling bool   `json:"reconsicling,omitempty" protobuf:"varint,1,opt,name=reconciling"`
	Etag        string `json:"etag,omitempty" protobuf:"varint,2,opt,name=etag"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type BorgJob struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec BorgJobSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status BorgJobStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BorgJobList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Server
	Items []BorgJob `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// BorgJobSpec
type BorgJobSpec struct {
	// +optional
	Environment string `json:"environment,omitempty" protobuf:"bytes,1,opt,name=environment"`
	Cell        string `json:"cell,omitempty" protobuf:"bytes,2,opt,name=cell"`
	TaskCount   int32  `json:"task_count,omitempty" protobuf:"bytes,3,opt,name=task_count"`
}

// BorgJobStatus
type BorgJobStatus struct {
	// +optional
	Reconciling bool   `json:"reconsicling,omitempty" protobuf:"varint,1,opt,name=reconciling"`
	Etag        string `json:"etag,omitempty" protobuf:"varint,2,opt,name=etag"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Environment struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec EnvironmentSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// +optional
	Status EnvironmentStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Server
	Items []Environment `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// EnvironmentSpec
type EnvironmentSpec struct {
	// +optional
	DataRealm          string `json:"data_realm,omitempty" protobuf:"bytes,1,opt,name=data_realm"`
	ReleaseStage       string `json:"release_stage,omitempty" protobuf:"bytes,2,opt,name=release_stage"`
	RuntimeEnvironment string `json:"runtime_environment,omitempty" protobuf:"bytes,3,opt,name=runtime_environment"`
}

// EvironmentStatus
type EnvironmentStatus struct {
	// +optional
	Reconciling bool   `json:"reconsicling,omitempty" protobuf:"varint,1,opt,name=reconciling"`
	Etag        string `json:"etag,omitempty" protobuf:"varint,2,opt,name=etag"`
}
