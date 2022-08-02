// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package azurearc

import (
	"context"
)

// StartClusterWatcher watches for CR
// Same as K8s
func (p *AzureArcProvider) StartClusterWatcher(ctx context.Context) error {
	return p.gitProvider.StartClusterWatcher(ctx)
}

// ApplyStatusCR applies status CR
func (p *AzureArcProvider) ApplyStatusCR(ctx context.Context, name string, content []byte) error {
	ref, err := p.gitProvider.Apply(ctx, name, nil, content)
	if err != nil {
		return err
	}
	return p.gitProvider.Commit(ctx, ref)
}

// DeleteStatusCR deletes status CR
func (p *AzureArcProvider) DeleteStatusCR(ctx context.Context, name string, content []byte) error {

	ref, err := p.gitProvider.Delete(name, nil, content)
	if err != nil {
		return err
	}
	return p.gitProvider.Commit(ctx, ref)
}
