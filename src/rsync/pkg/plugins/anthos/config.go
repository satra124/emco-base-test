// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package anthos

import (
	"context"
	"time"

	"github.com/fluxcd/go-git-providers/gitprovider"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	emcogit "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogit"
)

func (p *AnthosProvider) ApplyConfig(ctx context.Context, config interface{}) error {

	// Add to the commit
	path := "clusters/" + p.gitProvider.Cluster + "/context/" + p.gitProvider.Cid + "/deployed"
	var gp interface{}
	gp = emcogit.Add(path, time.Now().String(), []gitprovider.CommitFile{}, p.gitProvider.GitType)
	appName := p.gitProvider.Cid + p.gitProvider.App

	// Commit
	err := emcogit.CommitFiles(ctx, p.gitProvider.Client, p.gitProvider.UserName, p.gitProvider.RepoName, p.gitProvider.Branch, "Commit for "+p.gitProvider.GetPath("context"), appName, gp, p.gitProvider.GitType)
	if err != nil {
		log.Error("ApplyConfig:: Commit files err", log.Fields{"err": err, "gp": gp})
	}
	return err
}

func (p *AnthosProvider) DeleteConfig(ctx context.Context, config interface{}) error {
	path := "clusters/" + p.gitProvider.Cluster + "/" + p.gitProvider.Cid + ".yaml"
	var gp interface{}
	gp = emcogit.Delete(path, []gitprovider.CommitFile{}, p.gitProvider.GitType)
	path = "clusters/" + p.gitProvider.Cluster + "/" + "kcust" + p.gitProvider.Cid + ".yaml"
	appName := p.gitProvider.Cid + p.gitProvider.App

	err := emcogit.CommitFiles(ctx, p.gitProvider.Client, p.gitProvider.UserName, p.gitProvider.RepoName, p.gitProvider.Branch, "Commit for "+p.gitProvider.GetPath("context"), appName, gp, p.gitProvider.GitType)
	if err != nil {
		log.Error("ApplyConfig:: Commit files err", log.Fields{"err": err, "gp": gp})
	}
	return err
}
