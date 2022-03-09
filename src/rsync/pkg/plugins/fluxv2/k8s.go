// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package fluxv2

import (
	"fmt"
	pkgerrors "github.com/pkg/errors"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"strings"

	"github.com/fluxcd/go-git-providers/gitprovider"
	emcogithub "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogithub"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/internal/utils"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/db"
)

// Connection is for a cluster
type Fluxv2Provider struct {
	cid         string
	cluster     string
	app         string
	namespace   string
	level       string
	githubToken string
	userName    string
	branch      string
	repoName    string
	url         string
	client      gitprovider.Client
}

func NewFluxv2Provider(cid, app, cluster, level, namespace string) (*Fluxv2Provider, error) {

	result := strings.SplitN(cluster, "+", 2)

	c, err := utils.GetGitOpsConfig(cluster, level, namespace)
	if err != nil {
		return nil, err
	}
	if c.Props.GitOpsType != "fluxcd" {
		log.Error("Invalid GitOps type:", log.Fields{})
		return nil, pkgerrors.Errorf("Invalid GitOps type: " + c.Props.GitOpsType)
	}
	// Read from database
	ccc := db.NewCloudConfigClient()
	refObject, err := ccc.GetClusterSyncObjects(result[0], c.Props.GitOpsReferenceObject)
	if err != nil {
		log.Error("Invalid refObject :", log.Fields{"refObj": c.Props.GitOpsReferenceObject, "error": err})
		return nil, err
	}

	kv := refObject.Spec.Kv
	var githubToken, branch, userName, repoName string

	for _, kvpair := range kv {
		log.Info("kvpair", log.Fields{"kvpair": kvpair})
		v, ok := kvpair["githubToken"]
		if ok {
			githubToken = fmt.Sprintf("%v", v)
			continue
		}
		v, ok = kvpair["repoName"]
		if ok {
			repoName = fmt.Sprintf("%v", v)
			continue
		}
		v, ok = kvpair["userName"]
		if ok {
			userName = fmt.Sprintf("%v", v)
			continue
		}
		v, ok = kvpair["branch"]
		if ok {
			branch = fmt.Sprintf("%v", v)
			continue
		}
	}
	if len(githubToken) <= 0 || len(branch) <= 0 || len(userName) <= 0 || len(repoName) <= 0 {
		log.Error("Missing information for Github", log.Fields{"token": githubToken, "branch": branch, "userName": userName, "repoName": repoName})
		return nil, pkgerrors.Errorf("Missing Information for Github")
	}
	p := Fluxv2Provider{
		cid:         cid,
		app:         app,
		cluster:     cluster,
		level:       level,
		namespace:   namespace,
		githubToken: githubToken,
		branch:      branch,
		userName:    userName,
		repoName:    repoName,
		url:         "https://github.com/" + userName + "/" + repoName,
	}
	client, err := emcogithub.CreateClient(githubToken)
	if err != nil {
		log.Error("Error getting github client", log.Fields{"err": err})
		return nil, err
	}
	p.client = client
	return &p, nil
}

func (p *Fluxv2Provider) CleanClientProvider() error {
	return nil
}
