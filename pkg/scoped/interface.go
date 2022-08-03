package scoped

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
)

type ScopedSharedInformerFactory interface {
	Start(stopCh <-chan struct{})
	ForResource(gvr schema.GroupVersionResource) informers.GenericInformer
	WaitForCacheSync(stopCh <-chan struct{}) map[schema.GroupVersionResource]bool
}
