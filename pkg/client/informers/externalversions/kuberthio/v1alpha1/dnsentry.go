/*
Copyright 2018 The Openshift Evangelists

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	kuberthio_v1alpha1 "github.com/frodehus/kuberth/pkg/apis/kuberthio/v1alpha1"
	versioned "github.com/frodehus/kuberth/pkg/client/clientset/versioned"
	internalinterfaces "github.com/frodehus/kuberth/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/frodehus/kuberth/pkg/client/listers/kuberthio/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// DnsEntryInformer provides access to a shared informer and lister for
// DnsEntries.
type DnsEntryInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.DnsEntryLister
}

type dnsEntryInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewDnsEntryInformer constructs a new informer for DnsEntry type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDnsEntryInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDnsEntryInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredDnsEntryInformer constructs a new informer for DnsEntry type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDnsEntryInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KuberthV1alpha1().DnsEntries(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KuberthV1alpha1().DnsEntries(namespace).Watch(options)
			},
		},
		&kuberthio_v1alpha1.DnsEntry{},
		resyncPeriod,
		indexers,
	)
}

func (f *dnsEntryInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDnsEntryInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *dnsEntryInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kuberthio_v1alpha1.DnsEntry{}, f.defaultInformer)
}

func (f *dnsEntryInformer) Lister() v1alpha1.DnsEntryLister {
	return v1alpha1.NewDnsEntryLister(f.Informer().GetIndexer())
}