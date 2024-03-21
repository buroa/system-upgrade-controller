// Code generated by codegen. DO NOT EDIT.

package v1

import (
	"context"
	"sync"
	"time"

	v1 "github.com/buroa/system-upgrade-controller/pkg/apis/upgrade.cattle.io/v1"
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

type PlanHandler func(string, *v1.Plan) (*v1.Plan, error)

type PlanController interface {
	generic.ControllerMeta
	PlanClient

	OnChange(ctx context.Context, name string, sync PlanHandler)
	OnRemove(ctx context.Context, name string, sync PlanHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() PlanCache
}

type PlanClient interface {
	Create(*v1.Plan) (*v1.Plan, error)
	Update(*v1.Plan) (*v1.Plan, error)
	UpdateStatus(*v1.Plan) (*v1.Plan, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Plan, error)
	List(namespace string, opts metav1.ListOptions) (*v1.PlanList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Plan, err error)
}

type PlanCache interface {
	Get(namespace, name string) (*v1.Plan, error)
	List(namespace string, selector labels.Selector) ([]*v1.Plan, error)

	AddIndexer(indexName string, indexer PlanIndexer)
	GetByIndex(indexName, key string) ([]*v1.Plan, error)
}

type PlanIndexer func(obj *v1.Plan) ([]string, error)

type planController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewPlanController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) PlanController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &planController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromPlanHandlerToHandler(sync PlanHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Plan
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Plan))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *planController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Plan))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdatePlanDeepCopyOnChange(client PlanClient, obj *v1.Plan, handler func(obj *v1.Plan) (*v1.Plan, error)) (*v1.Plan, error) {
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

func (c *planController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *planController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *planController) OnChange(ctx context.Context, name string, sync PlanHandler) {
	c.AddGenericHandler(ctx, name, FromPlanHandlerToHandler(sync))
}

func (c *planController) OnRemove(ctx context.Context, name string, sync PlanHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromPlanHandlerToHandler(sync)))
}

func (c *planController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *planController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *planController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *planController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *planController) Cache() PlanCache {
	return &planCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *planController) Create(obj *v1.Plan) (*v1.Plan, error) {
	result := &v1.Plan{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *planController) Update(obj *v1.Plan) (*v1.Plan, error) {
	result := &v1.Plan{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *planController) UpdateStatus(obj *v1.Plan) (*v1.Plan, error) {
	result := &v1.Plan{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *planController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *planController) Get(namespace, name string, options metav1.GetOptions) (*v1.Plan, error) {
	result := &v1.Plan{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *planController) List(namespace string, opts metav1.ListOptions) (*v1.PlanList, error) {
	result := &v1.PlanList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *planController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *planController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1.Plan, error) {
	result := &v1.Plan{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type planCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *planCache) Get(namespace, name string) (*v1.Plan, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1.Plan), nil
}

func (c *planCache) List(namespace string, selector labels.Selector) (ret []*v1.Plan, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Plan))
	})

	return ret, err
}

func (c *planCache) AddIndexer(indexName string, indexer PlanIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Plan))
		},
	}))
}

func (c *planCache) GetByIndex(indexName, key string) (result []*v1.Plan, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.Plan, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.Plan))
	}
	return result, nil
}

// PlanStatusHandler is executed for every added or modified Plan. Should return the new status to be updated
type PlanStatusHandler func(obj *v1.Plan, status v1.PlanStatus) (v1.PlanStatus, error)

// PlanGeneratingHandler is the top-level handler that is executed for every Plan event. It extends PlanStatusHandler by a returning a slice of child objects to be passed to apply.Apply
type PlanGeneratingHandler func(obj *v1.Plan, status v1.PlanStatus) ([]runtime.Object, v1.PlanStatus, error)

// RegisterPlanStatusHandler configures a PlanController to execute a PlanStatusHandler for every events observed.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterPlanStatusHandler(ctx context.Context, controller PlanController, condition condition.Cond, name string, handler PlanStatusHandler) {
	statusHandler := &planStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromPlanHandlerToHandler(statusHandler.sync))
}

// RegisterPlanGeneratingHandler configures a PlanController to execute a PlanGeneratingHandler for every events observed, passing the returned objects to the provided apply.Apply.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterPlanGeneratingHandler(ctx context.Context, controller PlanController, apply apply.Apply,
	condition condition.Cond, name string, handler PlanGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &planGeneratingHandler{
		PlanGeneratingHandler: handler,
		apply:                 apply,
		name:                  name,
		gvk:                   controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterPlanStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type planStatusHandler struct {
	client    PlanClient
	condition condition.Cond
	handler   PlanStatusHandler
}

// sync is executed on every resource addition or modification. Executes the configured handlers and sends the updated status to the Kubernetes API
func (a *planStatusHandler) sync(key string, obj *v1.Plan) (*v1.Plan, error) {
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

type planGeneratingHandler struct {
	PlanGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
	seen  sync.Map
}

// Remove handles the observed deletion of a resource, cascade deleting every associated resource previously applied
func (a *planGeneratingHandler) Remove(key string, obj *v1.Plan) (*v1.Plan, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.Plan{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	if a.opts.UniqueApplyForResourceVersion {
		a.seen.Delete(key)
	}

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

// Handle executes the configured PlanGeneratingHandler and pass the resulting objects to apply.Apply, finally returning the new status of the resource
func (a *planGeneratingHandler) Handle(obj *v1.Plan, status v1.PlanStatus) (v1.PlanStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.PlanGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}
	if !a.isNewResourceVersion(obj) {
		return newStatus, nil
	}

	err = generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
	if err != nil {
		return newStatus, err
	}
	a.storeResourceVersion(obj)
	return newStatus, nil
}

// isNewResourceVersion detects if a specific resource version was already successfully processed.
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *planGeneratingHandler) isNewResourceVersion(obj *v1.Plan) bool {
	if !a.opts.UniqueApplyForResourceVersion {
		return true
	}

	// Apply once per resource version
	key := obj.Namespace + "/" + obj.Name
	previous, ok := a.seen.Load(key)
	return !ok || previous != obj.ResourceVersion
}

// storeResourceVersion keeps track of the latest resource version of an object for which Apply was executed
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *planGeneratingHandler) storeResourceVersion(obj *v1.Plan) {
	if !a.opts.UniqueApplyForResourceVersion {
		return
	}

	key := obj.Namespace + "/" + obj.Name
	a.seen.Store(key, obj.ResourceVersion)
}
