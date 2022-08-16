//=======================================================================
// Copyright (c) 2022 Aarna Networks, Inc.
// All rights reserved.
// ======================================================================
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//           http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ========================================================================

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "edgemetricscollector/pkg/apis/edgemetricscollector/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MetricsCollectorLister helps list MetricsCollectors.
// All objects returned here must be treated as read-only.
type MetricsCollectorLister interface {
	// List lists all MetricsCollectors in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.MetricsCollector, err error)
	// MetricsCollectors returns an object that can list and get MetricsCollectors.
	MetricsCollectors(namespace string) MetricsCollectorNamespaceLister
	MetricsCollectorListerExpansion
}

// metricsCollectorLister implements the MetricsCollectorLister interface.
type metricsCollectorLister struct {
	indexer cache.Indexer
}

// NewMetricsCollectorLister returns a new MetricsCollectorLister.
func NewMetricsCollectorLister(indexer cache.Indexer) MetricsCollectorLister {
	return &metricsCollectorLister{indexer: indexer}
}

// List lists all MetricsCollectors in the indexer.
func (s *metricsCollectorLister) List(selector labels.Selector) (ret []*v1alpha1.MetricsCollector, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.MetricsCollector))
	})
	return ret, err
}

// MetricsCollectors returns an object that can list and get MetricsCollectors.
func (s *metricsCollectorLister) MetricsCollectors(namespace string) MetricsCollectorNamespaceLister {
	return metricsCollectorNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// MetricsCollectorNamespaceLister helps list and get MetricsCollectors.
// All objects returned here must be treated as read-only.
type MetricsCollectorNamespaceLister interface {
	// List lists all MetricsCollectors in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.MetricsCollector, err error)
	// Get retrieves the MetricsCollector from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.MetricsCollector, error)
	MetricsCollectorNamespaceListerExpansion
}

// metricsCollectorNamespaceLister implements the MetricsCollectorNamespaceLister
// interface.
type metricsCollectorNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all MetricsCollectors in the indexer for a given namespace.
func (s metricsCollectorNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.MetricsCollector, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.MetricsCollector))
	})
	return ret, err
}

// Get retrieves the MetricsCollector from the indexer for a given namespace and name.
func (s metricsCollectorNamespaceLister) Get(name string) (*v1alpha1.MetricsCollector, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("metricscollector"), name)
	}
	return obj.(*v1alpha1.MetricsCollector), nil
}
