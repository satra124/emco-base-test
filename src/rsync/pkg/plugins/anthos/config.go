// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package anthos

import (
	"context"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"

)

func (p *AnthosProvider) ApplyConfig(ctx context.Context, config interface{}) error {

	// Add to the commit
	path := "clusters/" + p.gitProvider.Cluster + "/context/" + p.gitProvider.Cid + "/deployed"
	var gp interface{}
	gp, err := p.gitProvider.Add(path, path, "", gp)
	//appName := p.gitProvider.Cid + p.gitProvider.App

	// Commit
	err = p.gitProvider.Commit(ctx, gp)
	if err != nil {
		log.Error("ApplyConfig:: Commit files err", log.Fields{"err": err, "gp": gp})
	}
	return err
}

func (p *AnthosProvider) DeleteConfig(ctx context.Context, config interface{}) error {
	path := "clusters/" + p.gitProvider.Cluster + "/" + p.gitProvider.Cid + ".yaml"
	var gp interface{}
	gp, err := p.gitProvider.Delete(path, gp, nil)
	path = "clusters/" + p.gitProvider.Cluster + "/" + "kcust" + p.gitProvider.Cid + ".yaml"
//	appName := p.gitProvider.Cid + p.gitProvider.App

	err = p.gitProvider.Commit(ctx, gp)
	if err != nil {
		log.Error("ApplyConfig:: Commit files err", log.Fields{"err": err, "gp": gp})
	}
	return err
}
