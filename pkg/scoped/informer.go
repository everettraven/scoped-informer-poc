package scoped

import (
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamiclister"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type scopedSharedInformerFactory struct {
	client        dynamic.Interface
	defaultResync time.Duration

	lock      sync.Mutex
	informers map[schema.GroupVersionResource]informers.GenericInformer

	startedInformers map[schema.GroupVersionResource]bool
	// TODO: add a way to handle TweakListOptions in the future.
	// For the time being focus on core implementation before allowing list tweaks

	// TODO: consider adding a way to specify the namespace in the future.
	// For the time being focus on all namespace based
}

var _ ScopedSharedInformerFactory = &scopedSharedInformerFactory{}

func NewScopedSharedInformerFactory(client dynamic.Interface, defaultResync time.Duration) ScopedSharedInformerFactory {
	return &scopedSharedInformerFactory{
		client:           client,
		defaultResync:    defaultResync,
		informers:        make(map[schema.GroupVersionResource]informers.GenericInformer),
		startedInformers: make(map[schema.GroupVersionResource]bool),
	}
}

func (si *scopedSharedInformerFactory) Start(stopCh <-chan struct{}) {
	// TODO: implement this function
}

func (si *scopedSharedInformerFactory) ForResource(gvr schema.GroupVersionResource) informers.GenericInformer {
	// TODO: implement this function
	return nil
}

func (si *scopedSharedInformerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[schema.GroupVersionResource]bool {
	// TODO: implement this function
	return nil
}

type scopedInformer struct {
	informer cache.SharedIndexInformer
	gvr      schema.GroupVersionResource
}

func NewScopedInformer(client dynamic.Interface, gvr schema.GroupVersionResource, resyncPeriod time.Duration, indexers cache.Indexers) informers.GenericInformer {
	return &scopedInformer{
		gvr: gvr,
		informer: cache.NewSharedIndexInformer(
			NewScopedListerWatcher(client, gvr),
			&unstructured.Unstructured{},
			resyncPeriod,
			indexers),
	}
}

var _ informers.GenericInformer = &scopedInformer{}

func (si *scopedInformer) Informer() cache.SharedIndexInformer {
	return si.informer
}

func (si *scopedInformer) Lister() cache.GenericLister {
	// TODO: implement this properly
	// Question: From what I can tell, the dynamic lister in client-go uses the SharedIndexInformer Indexers for listing resources.
	// Is the Indexer (cache?) populated via the ListerWatcher that the SharedIndexInformer is created with?
	return dynamiclister.NewRuntimeObjectShim(dynamiclister.New(si.informer.GetIndexer(), si.gvr))
}
