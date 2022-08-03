package scoped

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// ScopedListerWatcherConfiguration is a struct meant to implement the
// ListerWatcherConfiguration interface
type ScopedListerWatcherConfiguration struct{}

func NewScopedListerWatcherConfiguration() *ScopedListerWatcherConfiguration {
	return &ScopedListerWatcherConfiguration{}
}

func (slwc *ScopedListerWatcherConfiguration) GetStructuredListerWatcher(config *rest.Config, mapping *meta.RESTMapping, scheme runtime.Scheme) (cache.ListerWatcher, error) {
	cli, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("encountered an error creating the client for the ScopedListerWatcher")
	}

	lw := NewScopedListerWatcher(cli, mapping.Resource)

	gvk := mapping.GroupVersionKind
	listGVK := gvk.GroupVersion().WithKind(gvk.Kind + "List")
	listObj, err := scheme.New(listGVK)
	if err != nil {
		return nil, err
	}

	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			l, err := lw.List(options)
			if err != nil {
				return nil, err
			}
			ul := l.(*unstructured.UnstructuredList)
			uc := ul.UnstructuredContent()
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(uc, listObj)
			if err != nil {
				return nil, err
			}

			return listObj, nil
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return lw.Watch(options)
		},
	}, nil
}

func (slwc *ScopedListerWatcherConfiguration) GetUnstructuredListerWatcher(config *rest.Config, mapping *meta.RESTMapping, scheme runtime.Scheme) (cache.ListerWatcher, error) {
	cli, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("encountered an error creating the client for the ScopedListerWatcher")
	}

	lw := NewScopedListerWatcher(cli, mapping.Resource)

	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return lw.List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return lw.Watch(options)
		},
	}, nil
}

func (slwc *ScopedListerWatcherConfiguration) GetMetadataListerWatcher(config *rest.Config, mapping *meta.RESTMapping, scheme runtime.Scheme) (cache.ListerWatcher, error) {
	cli, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("encountered an error creating the client for the ScopedListerWatcher")
	}

	lw := NewScopedListerWatcher(cli, mapping.Resource)

	gvk := mapping.GroupVersionKind
	list := &metav1.PartialObjectMetadataList{}

	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			l, err := lw.List(options)
			if err != nil {
				return nil, err
			}
			ul := l.(*unstructured.UnstructuredList)
			uc := ul.UnstructuredContent()
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(uc, list)
			if err != nil {
				return nil, err
			}

			for i := range list.Items {
				list.Items[i].SetGroupVersionKind(gvk)
			}

			return list, nil
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			watcher, err := lw.Watch(options)
			if err != nil {
				return nil, err
			}

			return newGVKFixupWatcher(gvk, watcher), nil
		},
	}, nil
}

func newGVKFixupWatcher(gvk schema.GroupVersionKind, watcher watch.Interface) watch.Interface {
	return watch.Filter(
		watcher,
		func(in watch.Event) (watch.Event, bool) {
			in.Object.GetObjectKind().SetGroupVersionKind(gvk)
			return in, true
		},
	)
}
