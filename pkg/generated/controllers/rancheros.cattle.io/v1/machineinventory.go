/*
Copyright 2022 SUSE LLC

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/rancher-sandbox/rancheros-operator/pkg/apis/rancheros.cattle.io/v1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type MachineInventoryHandler func(string, *v1.MachineInventory) (*v1.MachineInventory, error)

type MachineInventoryController interface {
	generic.ControllerMeta
	MachineInventoryClient

	OnChange(ctx context.Context, name string, sync MachineInventoryHandler)
	OnRemove(ctx context.Context, name string, sync MachineInventoryHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() MachineInventoryCache
}

type MachineInventoryClient interface {
	Create(*v1.MachineInventory) (*v1.MachineInventory, error)
	Update(*v1.MachineInventory) (*v1.MachineInventory, error)
	UpdateStatus(*v1.MachineInventory) (*v1.MachineInventory, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.MachineInventory, error)
	List(namespace string, opts metav1.ListOptions) (*v1.MachineInventoryList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.MachineInventory, err error)
}

type MachineInventoryCache interface {
	Get(namespace, name string) (*v1.MachineInventory, error)
	List(namespace string, selector labels.Selector) ([]*v1.MachineInventory, error)

	AddIndexer(indexName string, indexer MachineInventoryIndexer)
	GetByIndex(indexName, key string) ([]*v1.MachineInventory, error)
}

type MachineInventoryIndexer func(obj *v1.MachineInventory) ([]string, error)

type machineInventoryController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewMachineInventoryController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) MachineInventoryController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &machineInventoryController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromMachineInventoryHandlerToHandler(sync MachineInventoryHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.MachineInventory
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.MachineInventory))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *machineInventoryController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.MachineInventory))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateMachineInventoryDeepCopyOnChange(client MachineInventoryClient, obj *v1.MachineInventory, handler func(obj *v1.MachineInventory) (*v1.MachineInventory, error)) (*v1.MachineInventory, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *machineInventoryController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *machineInventoryController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *machineInventoryController) OnChange(ctx context.Context, name string, sync MachineInventoryHandler) {
	c.AddGenericHandler(ctx, name, FromMachineInventoryHandlerToHandler(sync))
}

func (c *machineInventoryController) OnRemove(ctx context.Context, name string, sync MachineInventoryHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromMachineInventoryHandlerToHandler(sync)))
}

func (c *machineInventoryController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *machineInventoryController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *machineInventoryController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *machineInventoryController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *machineInventoryController) Cache() MachineInventoryCache {
	return &machineInventoryCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *machineInventoryController) Create(obj *v1.MachineInventory) (*v1.MachineInventory, error) {
	result := &v1.MachineInventory{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *machineInventoryController) Update(obj *v1.MachineInventory) (*v1.MachineInventory, error) {
	result := &v1.MachineInventory{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *machineInventoryController) UpdateStatus(obj *v1.MachineInventory) (*v1.MachineInventory, error) {
	result := &v1.MachineInventory{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *machineInventoryController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *machineInventoryController) Get(namespace, name string, options metav1.GetOptions) (*v1.MachineInventory, error) {
	result := &v1.MachineInventory{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *machineInventoryController) List(namespace string, opts metav1.ListOptions) (*v1.MachineInventoryList, error) {
	result := &v1.MachineInventoryList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *machineInventoryController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *machineInventoryController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.MachineInventory, error) {
	result := &v1.MachineInventory{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type machineInventoryCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *machineInventoryCache) Get(namespace, name string) (*v1.MachineInventory, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.MachineInventory), nil
}

func (c *machineInventoryCache) List(namespace string, selector labels.Selector) (ret []*v1.MachineInventory, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.MachineInventory))
	})

	return ret, err
}

func (c *machineInventoryCache) AddIndexer(indexName string, indexer MachineInventoryIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.MachineInventory))
		},
	}))
}

func (c *machineInventoryCache) GetByIndex(indexName, key string) (result []*v1.MachineInventory, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.MachineInventory, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.MachineInventory))
	}
	return result, nil
}

type MachineInventoryStatusHandler func(obj *v1.MachineInventory, status v1.MachineInventoryStatus) (v1.MachineInventoryStatus, error)

type MachineInventoryGeneratingHandler func(obj *v1.MachineInventory, status v1.MachineInventoryStatus) ([]runtime.Object, v1.MachineInventoryStatus, error)

func RegisterMachineInventoryStatusHandler(ctx context.Context, controller MachineInventoryController, condition condition.Cond, name string, handler MachineInventoryStatusHandler) {
	statusHandler := &machineInventoryStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromMachineInventoryHandlerToHandler(statusHandler.sync))
}

func RegisterMachineInventoryGeneratingHandler(ctx context.Context, controller MachineInventoryController, apply apply.Apply,
	condition condition.Cond, name string, handler MachineInventoryGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &machineInventoryGeneratingHandler{
		MachineInventoryGeneratingHandler: handler,
		apply:                             apply,
		name:                              name,
		gvk:                               controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterMachineInventoryStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type machineInventoryStatusHandler struct {
	client    MachineInventoryClient
	condition condition.Cond
	handler   MachineInventoryStatusHandler
}

func (a *machineInventoryStatusHandler) sync(key string, obj *v1.MachineInventory) (*v1.MachineInventory, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type machineInventoryGeneratingHandler struct {
	MachineInventoryGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *machineInventoryGeneratingHandler) Remove(key string, obj *v1.MachineInventory) (*v1.MachineInventory, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.MachineInventory{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *machineInventoryGeneratingHandler) Handle(obj *v1.MachineInventory, status v1.MachineInventoryStatus) (v1.MachineInventoryStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.MachineInventoryGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
