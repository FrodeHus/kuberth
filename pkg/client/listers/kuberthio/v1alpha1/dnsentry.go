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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/frodehus/kuberth/pkg/apis/kuberthio/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DnsEntryLister helps list DnsEntries.
type DnsEntryLister interface {
	// List lists all DnsEntries in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.DnsEntry, err error)
	// DnsEntries returns an object that can list and get DnsEntries.
	DnsEntries(namespace string) DnsEntryNamespaceLister
	DnsEntryListerExpansion
}

// dnsEntryLister implements the DnsEntryLister interface.
type dnsEntryLister struct {
	indexer cache.Indexer
}

// NewDnsEntryLister returns a new DnsEntryLister.
func NewDnsEntryLister(indexer cache.Indexer) DnsEntryLister {
	return &dnsEntryLister{indexer: indexer}
}

// List lists all DnsEntries in the indexer.
func (s *dnsEntryLister) List(selector labels.Selector) (ret []*v1alpha1.DnsEntry, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DnsEntry))
	})
	return ret, err
}

// DnsEntries returns an object that can list and get DnsEntries.
func (s *dnsEntryLister) DnsEntries(namespace string) DnsEntryNamespaceLister {
	return dnsEntryNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DnsEntryNamespaceLister helps list and get DnsEntries.
type DnsEntryNamespaceLister interface {
	// List lists all DnsEntries in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.DnsEntry, err error)
	// Get retrieves the DnsEntry from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.DnsEntry, error)
	DnsEntryNamespaceListerExpansion
}

// dnsEntryNamespaceLister implements the DnsEntryNamespaceLister
// interface.
type dnsEntryNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all DnsEntries in the indexer for a given namespace.
func (s dnsEntryNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.DnsEntry, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DnsEntry))
	})
	return ret, err
}

// Get retrieves the DnsEntry from the indexer for a given namespace and name.
func (s dnsEntryNamespaceLister) Get(name string) (*v1alpha1.DnsEntry, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("dnsentry"), name)
	}
	return obj.(*v1alpha1.DnsEntry), nil
}
