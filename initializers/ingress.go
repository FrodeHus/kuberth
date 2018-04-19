package initializers

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/frodehus/kuberth/azure"
	"k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	annotation      = "initializer.kuberth.io/ingress"
	initializerName = "ingress.initializer.kuberth.io"
)

//Start the ingress initializer - Updating DNS server with configured ingress host
func Start(clientset *kubernetes.Clientset, stopCh <-chan struct{}) {
	glog.Info("Starting ingress initializer...")
	restClient := clientset.ExtensionsV1beta1().RESTClient()
	watchlist := cache.NewListWatchFromClient(restClient, "ingresses", corev1.NamespaceAll, fields.Everything())
	includeUninitializedWatchlist := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.IncludeUninitialized = true
			return watchlist.List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.IncludeUninitialized = true
			return watchlist.Watch(options)
		},
	}

	resyncPeriod := 30 * time.Second

	_, controller := cache.NewInformer(includeUninitializedWatchlist, &extv1beta1.Ingress{}, resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				err := initializeIngress(obj.(*extv1beta1.Ingress), clientset)
				if err != nil {
					log.Println(err)
				}
			},
		},
	)
	go controller.Run(stopCh)
}

func initializeIngress(ingress *extv1beta1.Ingress, clientset *kubernetes.Clientset) error {
	if ingress.ObjectMeta.GetInitializers() != nil {
		pendingInitializers := ingress.ObjectMeta.GetInitializers().Pending

		if initializerName == pendingInitializers[0].Name {
			log.Printf("Initializing ingress: %s", ingress.Name)

			initializedIngress := ingress.DeepCopy()

			// Remove self from the list of pending Initializers while preserving ordering.
			if len(pendingInitializers) == 1 {
				initializedIngress.ObjectMeta.Initializers = nil
			} else {
				initializedIngress.ObjectMeta.Initializers.Pending = append(pendingInitializers[:0], pendingInitializers[1:]...)
			}
			if ok := shouldInitialize(ingress, clientset, initializedIngress); !ok {
				_, err := clientset.ExtensionsV1beta1().Ingresses(ingress.Namespace).Update(initializedIngress)
				if err != nil {
					return err
				}
				return nil
			}

			createOrUpdateRecord(ingress)

			oldData, err := json.Marshal(ingress)
			if err != nil {
				return err
			}

			newData, err := json.Marshal(initializedIngress)
			if err != nil {
				return err
			}

			patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1beta1.Deployment{})
			if err != nil {
				return err
			}

			_, err = clientset.ExtensionsV1beta1().Ingresses(ingress.Namespace).Patch(ingress.Name, types.StrategicMergePatchType, patchBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func shouldInitialize(ingress *extv1beta1.Ingress, clientset *kubernetes.Clientset, initializedIngress *extv1beta1.Ingress) bool {
	a := ingress.ObjectMeta.GetAnnotations()
	_, ok := a[annotation]
	if !ok {
		log.Printf("Required '%s' annotation missing; skipping envoy container injection", annotation)
		return false
	}
	return true
}

func createOrUpdateRecord(ingress *extv1beta1.Ingress) (bool, error) {
	dnsProvider, err := azure.NewDNSClient()
	if err != nil {
		glog.Errorf("failed creating DNS provider: %s", err.Error())
		return false, err
	}
	for _, rule := range ingress.Spec.Rules {
		glog.Infof("Creating record for %s", rule.Host)
		record := rule.Host[0:strings.Index(rule.Host, ".")]
		_, err = dnsProvider.LookupRecord(record)
		if err != nil {
			glog.Errorf("failed creating A record: %s", err.Error())
			return false, err
		}
	}
	return true, nil
}
