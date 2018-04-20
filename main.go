package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"

	initializers "github.com/frodehus/kuberth/initializers"
	kuberth "github.com/frodehus/kuberth/pkg/client/clientset/versioned"
	informers "github.com/frodehus/kuberth/pkg/client/informers/externalversions"
	"github.com/frodehus/kuberth/pkg/signals"
	"github.com/golang/glog"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	flag.Parse()
}
func main() {
	glog.Info("Waking up Kuberth...")
	stopCh := signals.SetupSignalHandler()
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Info("Failed to init using kubeconfig, assuming in-cluster...")
		config, err = rest.InClusterConfig()
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	go initializers.Start(clientset, stopCh)
	dnsClient, err := kuberth.NewForConfig(config)
	informerFactory := informers.NewSharedInformerFactory(dnsClient, time.Second*30)
	go informerFactory.Start(stopCh)
	controller := NewDNSController(clientset, dnsClient, informerFactory)
	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}
