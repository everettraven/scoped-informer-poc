package main

import (
	"fmt"
	"os"

	"github.com/everettraven/scoped-informer-poc/pkg/scoped"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// For now this just serves as an example of how the scoped informer can be used by implementing it in a very simple Go program

func main() {
	UnstructuredTest()
	// StructuredTest()
	// MetadataTest()
}

func UnstructuredTest() {
	fmt.Println("Getting Pods!")
	cfg := config.GetConfigOrDie()
	gvk := corev1.SchemeGroupVersion.WithKind("Pod")
	mapper, err := apiutil.NewDiscoveryRESTMapper(cfg)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(1)
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(2)
	}

	scheme := runtime.NewScheme()

	lwCfg := scoped.NewScopedListerWatcherConfiguration()
	lw, err := lwCfg.GetUnstructuredListerWatcher(config.GetConfigOrDie(), mapping, scheme)
	pods, err := lw.List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(3)
	}

	fmt.Println("Got some Pods!")
	ulpods := pods.(*unstructured.UnstructuredList)

	for _, pod := range ulpods.Items {
		fmt.Println("Got Pod --> ", pod.GetName())
	}

	fmt.Println("-----------------------------------")
	fmt.Println("Watching Pods!")
	wc := make(chan string)
	watcher, err := lw.Watch(metav1.ListOptions{})
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(4)
	}

	go func() {
		for event := range watcher.ResultChan() {
			metaObj := event.Object.(metav1.Object)
			wc <- fmt.Sprintf("%s - `%s` in namespace `%s`", event.Type, metaObj.GetName(), metaObj.GetNamespace())
		}
	}()

	for {
		out := <-wc
		fmt.Println(out)
	}
}

func StructuredTest() {
	fmt.Println("Getting Pods!")
	cfg := config.GetConfigOrDie()
	gvk := corev1.SchemeGroupVersion.WithKind("Pod")
	mapper, err := apiutil.NewDiscoveryRESTMapper(cfg)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(1)
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(2)
	}

	scheme := runtime.NewScheme()

	err = corev1.AddToScheme(scheme)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(3)
	}

	lwCfg := scoped.NewScopedListerWatcherConfiguration()
	lw, err := lwCfg.GetStructuredListerWatcher(config.GetConfigOrDie(), mapping, scheme)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(4)
	}
	pods, err := lw.List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(5)
	}

	fmt.Println("Got some Pods!")
	podList := pods.(*corev1.PodList)

	for _, pod := range podList.Items {
		fmt.Println("Got Pod --> ", pod.GetName())
	}

	fmt.Println("-----------------------------------")
	fmt.Println("Watching Pods!")
	wc := make(chan string)
	watcher, err := lw.Watch(metav1.ListOptions{})
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(6)
	}

	go func() {
		for event := range watcher.ResultChan() {
			pod := event.Object.(*corev1.Pod)
			wc <- fmt.Sprintf("%s - `%s` in namespace `%s`", event.Type, pod.GetName(), pod.GetNamespace())
		}
	}()

	for {
		out := <-wc
		fmt.Println(out)
	}
}

func MetadataTest() {
	fmt.Println("Getting Pods!")
	cfg := config.GetConfigOrDie()
	gvk := corev1.SchemeGroupVersion.WithKind("Pod")
	mapper, err := apiutil.NewDiscoveryRESTMapper(cfg)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(1)
	}

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(2)
	}

	scheme := runtime.NewScheme()

	err = corev1.AddToScheme(scheme)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(3)
	}

	lwCfg := scoped.NewScopedListerWatcherConfiguration()
	lw, err := lwCfg.GetMetadataListerWatcher(config.GetConfigOrDie(), mapping, scheme)
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(4)
	}
	pods, err := lw.List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(5)
	}

	fmt.Println("Got some Pods!")
	mdList := pods.(*metav1.PartialObjectMetadataList)

	for _, pod := range mdList.Items {
		fmt.Println("Got Pod --> ", pod.GetName())
	}

	fmt.Println("-----------------------------------")
	fmt.Println("Watching Pods!")
	wc := make(chan string)
	watcher, err := lw.Watch(metav1.ListOptions{})
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(6)
	}

	go func() {
		for event := range watcher.ResultChan() {
			pod := event.Object.(*metav1.PartialObjectMetadata)
			wc <- fmt.Sprintf("%s - `%s` in namespace `%s`", event.Type, pod.GetName(), pod.GetNamespace())
		}
	}()

	for {
		out := <-wc
		fmt.Println(out)
	}
}
