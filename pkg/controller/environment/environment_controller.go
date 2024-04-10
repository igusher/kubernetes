/*
Copyright 2015 The Kubernetes Authors.
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

// Package server contains all the logic for handling Kubernetes Deployments.
// It implements a set of strategies (rolling, recreate) for deploying an application,
// the means to rollback to previous versions, proportional scaling for mitigating
// risk, cleanup policy, and other useful features of Deployments.
package environment

import (
	"context"
	"fmt"
	"time"

	"k8s.io/klog/v2"

	serverplatformv1 "k8s.io/api/serverplatform/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	serverplatforminformers "k8s.io/client-go/informers/serverplatform/v1"
	clientset "k8s.io/client-go/kubernetes"
	serverlisters "k8s.io/client-go/listers/serverplatform/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/controller"
)

// controllerKind contains the schema.GroupVersionKind for this controller type.
var controllerKind = serverplatformv1.SchemeGroupVersion.WithKind("Environment")

// DeploymentController is responsible for synchronizing Deployment objects stored
// in the system with actual running replica sets and pods.
type EnvironmentController struct {
	client    clientset.Interface

	// To allow injection of syncServer for testing.
	syncHandler func(ctx context.Context, dKey string) error
	// used for unit testing
	enqueueEnvironment func(server *serverplatformv1.Environment)

	// envLister can list/get servers from the shared informer's store
	envLister serverlisters.EnvironmentLister

	// Added as a member to the struct to allow injection for testing.
	envListerSynced cache.InformerSynced

	// Deployments that need to be synced
	queue workqueue.RateLimitingInterface
}

// NewDeploymentController creates a new DeploymentController.
func NewEnvironmentController(envInformer serverplatforminformers.EnvironmentInformer, client clientset.Interface) (*EnvironmentController, error) {
	dc := &EnvironmentController{
		client:           client,
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "serverplatform"),
	}

	envInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    dc.addEnvironment,
		UpdateFunc: dc.updateEnvironment,
		// This will enter the sync loop and no-op, because the server has been deleted from the store.
		DeleteFunc: dc.deleteEnvironment,
	})

	dc.syncHandler = dc.syncEnvironment
	dc.enqueueEnvironment = dc.enqueue

	dc.envLister = envInformer.Lister()
	dc.envListerSynced = envInformer.Informer().HasSynced
	return dc, nil
}

// Run begins watching and syncing.
func (dc *EnvironmentController) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()
	defer dc.queue.ShutDown()

	klog.InfoS("Starting controller", "environment", "serverplatform")
	defer klog.InfoS("Shutting down controller", "environment", "serverplatform")

	if !cache.WaitForNamedCacheSync("serverplatform", ctx.Done(), dc.envListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, dc.worker, time.Second)
	}

	<-ctx.Done()
}

func (dc *EnvironmentController) addEnvironment(obj interface{}) {
	d := obj.(*serverplatformv1.Environment)
	klog.InfoS("Adding environment", "serverplatform", klog.KObj(d))
	dc.enqueueEnvironment(d)
}

func (dc *EnvironmentController) updateEnvironment(old, cur interface{}) {
	oldD := old.(*serverplatformv1.Environment)
	curD := cur.(*serverplatformv1.Environment)
	klog.InfoS("Updating environment", "serverplatform", klog.KObj(oldD))
	dc.enqueueEnvironment(curD)
}

func (dc *EnvironmentController) deleteEnvironment(obj interface{}) {
	env, _ := obj.(*serverplatformv1.Environment)
	klog.InfoS("Deleting server", "serverplatform", klog.KObj(env))
	dc.enqueueEnvironment(env)
}


func (dc *EnvironmentController) enqueue(serverplatform *serverplatformv1.Environment) {
	key, err := controller.KeyFunc(serverplatform)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", serverplatform, err))
		return
	}

	dc.queue.Add(key)
}

// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (dc *EnvironmentController) worker(ctx context.Context) {
	for dc.processNextWorkItem(ctx) {
	}
}

func (dc *EnvironmentController) processNextWorkItem(ctx context.Context) bool {
	key, quit := dc.queue.Get()
	if quit {
		return false
	}
	defer dc.queue.Done(key)

	err := dc.syncHandler(ctx, key.(string))
	dc.handleErr(err, key)

	return true
}

func (dc *EnvironmentController) handleErr(err error, key interface{}) {
	if err == nil || errors.HasStatusCause(err, v1.NamespaceTerminatingCause) {
		dc.queue.Forget(key)
		return
	}
	utilruntime.HandleError(err)
	klog.V(2).InfoS("Dropping environment out of the queue", "serverplatform", key, "err", err)
	dc.queue.Forget(key)
}

// syncServer will sync the server with the given key.
// This function is not meant to be invoked concurrently with the same key.
func (dc *EnvironmentController) syncEnvironment(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", key)
		return err
	}

	startTime := time.Now()
	klog.V(4).InfoS("Started syncing environemnts", "env", klog.KRef(namespace, name), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing environments", "env", klog.KRef(namespace, name), "duration", time.Since(startTime))
	}()

	env, err := dc.envLister.Environments(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("Env has been deleted", "serverplatform", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return err
	}
	klog.V(2).InfoS("Env is found", "environment", klog.KObj(env))
	defaultCells := []string{"yc", "ab", "jk"}
	defaultTaskCount := 2
	for _, c := range defaultCells {
		borgjob := &serverplatformv1.BorgJob{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("borgjob-%s-%s", env.ObjectMeta.Name, c),
			},
			Spec: serverplatformv1.BorgJobSpec{
				Environment: env.ObjectMeta.Name,
				Cell: c,
				TaskCount: int32(defaultTaskCount),
			},
		}
		newBorgJob, err := dc.client.ServerplatformV1().BorgJobs(namespace).Create(ctx, borgjob, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		klog.V(2).InfoS("Create new BorgJob", "borgjob", klog.KObj(newBorgJob))
	}
	return nil
}
