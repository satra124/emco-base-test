// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package azurearc

import (
	"context"
)

// StartClusterWatcher watches for CR
// Same as K8s
func (p *AzureArcProvider) StartClusterWatcher() error {
	return p.gitProvider.StartClusterWatcher()
}

// ApplyStatusCR applies status CR
func (p *AzureArcProvider) ApplyStatusCR(name string, content []byte) error {
	ref, err := p.gitProvider.Apply(name, nil, content)
	if err != nil {
		return err
	}
	return p.gitProvider.Commit(context.Background(), ref)
}

// DeleteStatusCR deletes status CR
func (p *AzureArcProvider) DeleteStatusCR(name string, content []byte) error {

	ref, err := p.gitProvider.Delete(name, nil, content)
	if err != nil {
		return err
	}
	return p.gitProvider.Commit(context.Background(), ref)
}
