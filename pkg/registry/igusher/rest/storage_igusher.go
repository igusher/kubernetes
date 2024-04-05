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

package rest

import (
	"fmt"
	//"errors"

	igusherv1 "k8s.io/api/igusher/v1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	//utilfeature "k8s.io/apiserver/pkg/util/feature"
	//"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/igusher"
	//"k8s.io/kubernetes/pkg/features"
	igusherstore "k8s.io/kubernetes/pkg/registry/igusher/igusher/storage"
)

type RESTStorageProvider struct{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {

	grv := legacyscheme.Scheme.PrioritizedVersionsAllGroups()
	for _, g := range grv {
		fmt.Printf("#### %q - %q\n", g.Group, g.Version)
	}
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(igusher.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo
	storageMap, err := p.v1Storage(apiResourceConfigSource, restOptionsGetter)
	if err != nil {
		return apiGroupInfo, err
	}
	if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[igusherv1.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, nil

}

func (p RESTStorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// tokenreviews
	if resource := "igushers"; apiResourceConfigSource.ResourceEnabled(igusherv1.SchemeGroupVersion.WithResource(resource)) {
		igusherStorage, _, err := igusherstore.NewREST(restOptionsGetter)
		if err != nil {
			return storage, err
		}
		fmt.Printf("%v", igusherStorage)
		storage[resource] = igusherStorage
		//storage[resource + "/status"] = igusherStatus
	}

	return storage, nil
}

func (p RESTStorageProvider) GroupName() string {
	return igusher.GroupName
}
