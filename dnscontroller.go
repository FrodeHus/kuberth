package main

import (
	"fmt"
	"strings"
	"time"

	kuberth "github.com/frodehus/kuberth/pkg/client/clientset/versioned"
	informers "github.com/frodehus/kuberth/pkg/client/informers/externalversions"
	listers "github.com/frodehus/kuberth/pkg/client/listers/kuberthio/v1alpha1"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "kuberth-dns-controller"

type DNSController struct {
	kubeclientset   kubernetes.Interface
	kuberth         kuberth.Interface
	kuberthInformer informers.SharedInformerFactory
	workqueue       workqueue.RateLimitingInterface
	dnsLister       listers.DnsEntryLister
	dnsSynced       cache.InformerSynced
	recorder        record.EventRecorder
}

func NewDNSController(kubeclientset kubernetes.Interface, kuberth kuberth.Interface, kuberthInformer informers.SharedInformerFactory) *DNSController {
	dnsInformer := kuberthInformer.Kuberth().V1alpha1().DnsEntries()
	// eventBroadcaster := record.NewBroadcaster()
	// eventBroadcaster.StartLogging(glog.Infof)
	// eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	// recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &DNSController{
		kubeclientset:   kubeclientset,
		kuberth:         kuberth,
		kuberthInformer: kuberthInformer,
		dnsLister:       dnsInformer.Lister(),
		dnsSynced:       dnsInformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "DnsEntries"),
		// recorder:        recorder,
	}

	dnsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueDnsEntry,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueDnsEntry(new)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("Deleted object")
		},
	})
	return controller
}

func (c *DNSController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()
	glog.Info("Starting Kuberth DNS controller")
	if ok := cache.WaitForCacheSync(stopCh, c.dnsSynced); !ok {
		return fmt.Errorf("failed waiting for cache to sync")
	}
	glog.Infof("Starting %d workers", threadiness)
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Workers started")
	<-stopCh
	glog.Info("Shutting down workers")
	return nil
}

func (c *DNSController) runWorker() {
	for c.processNextWorkItem() {

	}
}

func (c *DNSController) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue, but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		c.workqueue.Forget(obj)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *DNSController) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	dnsEntry, err := c.dnsLister.DnsEntries(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	entries := dnsEntry.Spec
	if entries == nil {
		runtime.HandleError(fmt.Errorf("%s: no entries specified", key))
		return nil
	}
	glog.Infof("New record set [%s]\n", dnsEntry.Name)
	for _, record := range entries {
		glog.Infof("|%-12s|%-6s|%-16s\n", record.Name, strings.ToUpper(record.Type), record.Value)
		record.CreateOrUpdateRecord()
		// c.recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	}
	return nil
}

func (c *DNSController) enqueueDnsEntry(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
