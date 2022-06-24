// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package clusterprovider

import (
	"gitlab.com/project-emco/core/emco-base/src/ca-certs/pkg/module"
)

// ClusterGroupManager exposes all the caCert clusterGroup functionalities
type ClusterGroupManager interface {
	CreateClusterGroup(cluster module.ClusterGroup, cert, clusterProvider string, failIfExists bool) (module.ClusterGroup, bool, error)
	DeleteClusterGroup(cert, cluster, clusterProvider string) error
	GetAllClusterGroups(cert, clusterProvider string) ([]module.ClusterGroup, error)
	GetClusterGroup(cert, cluster, clusterProvider string) (module.ClusterGroup, error)
}

// ClusterGroupKey represents the resources associated with a caCert clusterGroup
type ClusterGroupKey struct {
	Cert            string `json:"caCertCp"`
	ClusterGroup    string `json:"caCertClusterGroupCp"`
	ClusterProvider string `json:"clusterProvider"`
}

// ClusterGroupClient holds the client properties
type ClusterGroupClient struct {
}

// NewClusterGroupClient returns an instance of the ClusterGroupClient which implements the Manager
func NewClusterGroupClient() *ClusterGroupClient {
	return &ClusterGroupClient{}
}

// CreateClusterGroup creates a caCert clusterGroup
func (c *ClusterGroupClient) CreateClusterGroup(group module.ClusterGroup, cert, clusterProvider string, failIfExists bool) (module.ClusterGroup, bool, error) {
	ck := ClusterGroupKey{
		Cert:            cert,
		ClusterGroup:    group.MetaData.Name,
		ClusterProvider: clusterProvider}

	return module.NewClusterGroupClient(ck).CreateClusterGroup(group, failIfExists)
}

// DeleteClusterGroup deletes a caCert clusterGroup
func (c *ClusterGroupClient) DeleteClusterGroup(cert, group, clusterProvider string) error {
	ck := ClusterGroupKey{
		Cert:            cert,
		ClusterGroup:    group,
		ClusterProvider: clusterProvider}

	return module.NewClusterGroupClient(ck).DeleteClusterGroup()
}

// GetAllClusterGroups returns all the caCert clusterGroup
func (c *ClusterGroupClient) GetAllClusterGroups(cert, clusterProvider string) ([]module.ClusterGroup, error) {
	ck := ClusterGroupKey{
		Cert:            cert,
		ClusterProvider: clusterProvider}

	return module.NewClusterGroupClient(ck).GetAllClusterGroups()
}

// GetClusterGroup returns the caCert clusterGroup
func (c *ClusterGroupClient) GetClusterGroup(cert, clusterGroup, clusterProvider string) (module.ClusterGroup, error) {
	ck := ClusterGroupKey{
		Cert:            cert,
		ClusterGroup:    clusterGroup,
		ClusterProvider: clusterProvider}

	return module.NewClusterGroupClient(ck).GetClusterGroup()
}
