// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package azurearcv2

import (
	"context"
)

// StartClusterWatcher watches for CR
// Same as K8s
func (p *AzureArcV2Provider) StartClusterWatcher(ctx context.Context) error {
	return p.gitProvider.StartClusterWatcher(ctx)
}

// ApplyStatusCR applies status CR
func (p *AzureArcV2Provider) ApplyStatusCR(ctx context.Context, name string, content []byte) error {
	ref, err := p.gitProvider.Apply(ctx, name, nil, content)
	if err != nil {
		return err
	}
	return p.gitProvider.Commit(ctx, ref)
}

// DeleteStatusCR deletes status CR
func (p *AzureArcV2Provider) DeleteStatusCR(ctx context.Context, name string, content []byte) error {

	ref, err := p.gitProvider.Delete(name, nil, content)
	if err != nil {
		return err
	}
	return p.gitProvider.Commit(ctx, ref)
}
