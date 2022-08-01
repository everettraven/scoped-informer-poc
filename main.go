package main

import (
	"fmt"
	"os"

	"github.com/everettraven/scoped-informer-poc/pkg/scoped"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// For now this just serves as an example of how the scoped informer can be used by implementing it in a very simple Go program

func main() {
	fmt.Println("Getting Pods!")
	cli := dynamic.NewForConfigOrDie(config.GetConfigOrDie())
	gvr := corev1.SchemeGroupVersion.WithResource("pods")
	lw := scoped.NewScopedListerWatcher(cli, gvr)
	pods, err := lw.List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("ERROR --> ", err)
		os.Exit(1)
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
		os.Exit(2)
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
