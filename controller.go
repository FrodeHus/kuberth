package main

import (
	kuberth "github.com/frodehus/kuberth/pkg/client/clientset/versioned"
	informers "github.com/frodehus/kuberth/pkg/client/informers/externalversions"
	listers "github.com/frodehus/kuberth/pkg/client/listers/kuberthio/v1alpha1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	kubeclientset   kubernetes.Interface
	kuberth         kuberth.Interface
	kuberthInformer informers.SharedInformerFactory
	workqueue       workqueue.RateLimitingInterface
	dnsLister       listers.DnsEntryLister
	dnsSynced       cache.InformerSynced
}

func NewController(kubeclientset kubernetes.Interface, kuberth kuberth.Interface, kuberthInformer informers.SharedInformerFactory) *Controller {
	dnsInformer := kuberthInformer.Kuberth().V1alpha1().DnsEntries()
	controller := &Controller{
		kubeclientset:   kubeclientset,
		kuberth:         kuberth,
		kuberthInformer: kuberthInformer,
		dnsLister:       dnsInformer.Lister(),
		dnsSynced:       dnsInformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "DnsEntries"),
	}
	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()
	return nil
}
