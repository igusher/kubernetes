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

// Package deployment contains all the logic for handling Kubernetes Deployments.
// It implements a set of strategies (rolling, recreate) for deploying an application,
// the means to rollback to previous versions, proportional scaling for mitigating
// risk, cleanup policy, and other useful features of Deployments.
package igusher

import (
	"context"
	"fmt"
	//"reflect"
	"time"

	"k8s.io/klog/v2"

	igusherv1 "k8s.io/api/igusher/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/labels"
	//"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	igusherinformers "k8s.io/client-go/informers/igusher/v1"
	//coreinformers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	igusherlisters "k8s.io/client-go/listers/igusher/v1"
	//corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/controller"
	//"k8s.io/kubernetes/pkg/controller/deployment/util"
)

// controllerKind contains the schema.GroupVersionKind for this controller type.
var controllerKind = igusherv1.SchemeGroupVersion.WithKind("Igusher")

// DeploymentController is responsible for synchronizing Deployment objects stored
// in the system with actual running replica sets and pods.
type IgusherController struct {
	client    clientset.Interface

	eventBroadcaster record.EventBroadcaster
	eventRecorder    record.EventRecorder

	// To allow injection of syncDeployment for testing.
	syncHandler func(ctx context.Context, dKey string) error
	// used for unit testing
	enqueueDeployment func(deployment *igusherv1.Igusher)

	// dLister can list/get deployments from the shared informer's store
	igLister igusherlisters.IgusherLister

	// dListerSynced returns true if the Deployment store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	igListerSynced cache.InformerSynced

	// Deployments that need to be synced
	queue workqueue.RateLimitingInterface
}

// NewDeploymentController creates a new DeploymentController.
func NewIgusherController(igInformer igusherinformers.IgusherInformer, client clientset.Interface) (*IgusherController, error) {
	eventBroadcaster := record.NewBroadcaster()

	dc := &IgusherController{
		client:           client,
		eventBroadcaster: eventBroadcaster,
		eventRecorder:    eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "igusher-controller"}),
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "igusher"),
	}

	igInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    dc.addIgusher,
		UpdateFunc: dc.updateIgusher,
		// This will enter the sync loop and no-op, because the deployment has been deleted from the store.
		DeleteFunc: dc.deleteIgusher,
	})

	dc.syncHandler = dc.syncDeployment
	dc.enqueueDeployment = dc.enqueue

	dc.igLister = igInformer.Lister()
	dc.igListerSynced = igInformer.Informer().HasSynced
	return dc, nil
}

// Run begins watching and syncing.
func (dc *IgusherController) Run(ctx context.Context, workers int) {
	defer utilruntime.HandleCrash()

	// Start events processing pipeline.
	dc.eventBroadcaster.StartStructuredLogging(0)
	dc.eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: dc.client.CoreV1().Events("")})
	defer dc.eventBroadcaster.Shutdown()

	defer dc.queue.ShutDown()

	klog.InfoS("Starting controller", "controller", "igusher")
	defer klog.InfoS("Shutting down controller", "controller", "igusher")

	if !cache.WaitForNamedCacheSync("igusher", ctx.Done(), dc.igListerSynced) {
		return
	}

	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, dc.worker, time.Second)
	}

	<-ctx.Done()
}

func (dc *IgusherController) addIgusher(obj interface{}) {
	d := obj.(*igusherv1.Igusher)
	klog.V(4).InfoS("Adding igusher", "igusher", klog.KObj(d))
	klog.InfoS("Adding igusher", "igusher", klog.KObj(d))
	dc.enqueueDeployment(d)
}

func (dc *IgusherController) updateIgusher(old, cur interface{}) {
	oldD := old.(*igusherv1.Igusher)
	curD := cur.(*igusherv1.Igusher)
	klog.V(4).InfoS("Updating igusher", "igusher", klog.KObj(oldD))
	klog.InfoS("Updating igusher", "igusher", klog.KObj(oldD))
	dc.enqueueDeployment(curD)
}

func (dc *IgusherController) deleteIgusher(obj interface{}) {
	//d, ok := obj.(*igusherv1.Igusher)
	d, _ := obj.(*igusherv1.Igusher)
	/*
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		d, ok = tombstone.Obj.(*apps.Deployment)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Deployment %#v", obj))
			return
		}
	}*/
	klog.V(4).InfoS("Deleting igusher", "igusher", klog.KObj(d))
	klog.InfoS("Deleting igusher", "igusher", klog.KObj(d))
	dc.enqueueDeployment(d)
}


func (dc *IgusherController) enqueue(igusher *igusherv1.Igusher) {
	key, err := controller.KeyFunc(igusher)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", igusher, err))
		return
	}

	dc.queue.Add(key)
}


/*
// resolveControllerRef returns the controller referenced by a ControllerRef,
// or nil if the ControllerRef could not be resolved to a matching controller
// of the correct Kind.
func (dc *IgusherController) resolveControllerRef(namespace string, controllerRef *metav1.OwnerReference) *apps.Deployment {
	// We can't look up by UID, so look up by Name and then verify UID.
	// Don't even try to look up by Name if it's the wrong Kind.
	if controllerRef.Kind != controllerKind.Kind {
		return nil
	}
	d, err := dc.igLister.Igushers(namespace).Get(controllerRef.Name)
	if err != nil {
		return nil
	}
	if d.UID != controllerRef.UID {
		// The controller we found with this Name is not the same one that the
		// ControllerRef points to.
		return nil
	}
	return d
}
*/
// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (dc *IgusherController) worker(ctx context.Context) {
	for dc.processNextWorkItem(ctx) {
	}
}

func (dc *IgusherController) processNextWorkItem(ctx context.Context) bool {
	key, quit := dc.queue.Get()
	if quit {
		return false
	}
	defer dc.queue.Done(key)

	err := dc.syncHandler(ctx, key.(string))
	dc.handleErr(err, key)

	return true
}

func (dc *IgusherController) handleErr(err error, key interface{}) {
	if err == nil || errors.HasStatusCause(err, v1.NamespaceTerminatingCause) {
		dc.queue.Forget(key)
		return
	}
	/*
	ns, name, keyErr := cache.SplitMetaNamespaceKey(key.(string))
	if keyErr != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", key)
	}

	if dc.queue.NumRequeues(key) < maxRetries {
		klog.V(2).InfoS("Error syncing deployment", "deployment", klog.KRef(ns, name), "err", err)
		dc.queue.AddRateLimited(key)
		return
	}
	*/
	utilruntime.HandleError(err)
	//klog.V(2).InfoS("Dropping igusher out of the queue", "igusher", klog.KRef(ns, name), "err", err)
	klog.V(2).InfoS("Dropping igusher out of the queue", "igusher", key, "err", err)
	dc.queue.Forget(key)
}

// syncDeployment will sync the deployment with the given key.
// This function is not meant to be invoked concurrently with the same key.
func (dc *IgusherController) syncDeployment(ctx context.Context, key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.ErrorS(err, "Failed to split meta namespace cache key", "cacheKey", key)
		return err
	}

	startTime := time.Now()
	klog.V(4).InfoS("Started syncing igusher", "igusher", klog.KRef(namespace, name), "startTime", startTime)
	defer func() {
		klog.V(4).InfoS("Finished syncing deployment", "deployment", klog.KRef(namespace, name), "duration", time.Since(startTime))
	}()
	
	deployment, err := dc.igLister.Igushers(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).InfoS("Igusher has been deleted", "igusher", klog.KRef(namespace, name))
		return nil
	}
	if err != nil {
		return err
	}
	klog.V(2).InfoS("Igusher is found", "igusher", klog.KObj(deployment))
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: deployment.Spec.Token,
		},
		Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name: deployment.Spec.Token,
						Image: "nginx:1.14.2",
					},
				},
		},
	}
	newPod, err := dc.client.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	klog.V(2).InfoS("Create new Pod", "pod", klog.KObj(newPod))
	/*
	// Deep-copy otherwise we are mutating our cache.
	// TODO: Deep-copy only when needed.
	d := deployment.DeepCopy()

	everything := metav1.LabelSelector{}
	if reflect.DeepEqual(d.Spec.Selector, &everything) {
		dc.eventRecorder.Eventf(d, v1.EventTypeWarning, "SelectingAll", "This deployment is selecting all pods. A non-empty selector is required.")
		if d.Status.ObservedGeneration < d.Generation {
			d.Status.ObservedGeneration = d.Generation
			dc.client.AppsV1().Deployments(d.Namespace).UpdateStatus(ctx, d, metav1.UpdateOptions{})
		}
		return nil
	}

	// List ReplicaSets owned by this Deployment, while reconciling ControllerRef
	// through adoption/orphaning.
	rsList, err := dc.getReplicaSetsForDeployment(ctx, d)
	if err != nil {
		return err
	}
	// List all Pods owned by this Deployment, grouped by their ReplicaSet.
	// Current uses of the podMap are:
	//
	// * check if a Pod is labeled correctly with the pod-template-hash label.
	// * check that no old Pods are running in the middle of Recreate Deployments.
	podMap, err := dc.getPodMapForDeployment(d, rsList)
	if err != nil {
		return err
	}

	if d.DeletionTimestamp != nil {
		return dc.syncStatusOnly(ctx, d, rsList)
	}

	// Update deployment conditions with an Unknown condition when pausing/resuming
	// a deployment. In this way, we can be sure that we won't timeout when a user
	// resumes a Deployment with a set progressDeadlineSeconds.
	if err = dc.checkPausedConditions(ctx, d); err != nil {
		return err
	}

	if d.Spec.Paused {
		return dc.sync(ctx, d, rsList)
	}

	// rollback is not re-entrant in case the underlying replica sets are updated with a new
	// revision so we should ensure that we won't proceed to update replica sets until we
	// make sure that the deployment has cleaned up its rollback spec in subsequent enqueues.
	if getRollbackTo(d) != nil {
		return dc.rollback(ctx, d, rsList)
	}

	scalingEvent, err := dc.isScalingEvent(ctx, d, rsList)
	if err != nil {
		return err
	}
	if scalingEvent {
		return dc.sync(ctx, d, rsList)
	}

	switch d.Spec.Strategy.Type {
	case apps.RecreateDeploymentStrategyType:
		return dc.rolloutRecreate(ctx, d, rsList, podMap)
	case apps.RollingUpdateDeploymentStrategyType:
		return dc.rolloutRolling(ctx, d, rsList)
	}
	return fmt.Errorf("unexpected deployment strategy type: %s", d.Spec.Strategy.Type)
	*/
	return nil
}
