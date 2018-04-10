package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"

	kuberth "github.com/frodehus/kuberth/pkg/client/clientset/versioned"
	informers "github.com/frodehus/kuberth/pkg/client/informers/externalversions"
	"github.com/frodehus/kuberth/pkg/signals"
	"github.com/golang/glog"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	stopCh := signals.SetupSignalHandler()
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	dnsClient, err := kuberth.NewForConfig(config)
	//kuberthAPI := dnsClient.KuberthV1alpha1()
	//listOptions := metav1.ListOptions{}
	//entries, err := kuberthAPI.DnsEntries("default").List(listOptions)
	informerFactory := informers.NewSharedInformerFactory(dnsClient, time.Second*30)
	go informerFactory.Start(stopCh)
	controller := NewController(clientset, dnsClient, informerFactory)
	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
