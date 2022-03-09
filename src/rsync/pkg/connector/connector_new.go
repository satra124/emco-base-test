// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package connector

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/plugins/fluxv2"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/plugins/k8s"
	. "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/types"
	"strings"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/internal/utils"
)

// Connection is for a cluster
type Provider struct {
	cid string
}

func NewProvider(id interface{}) Provider {
	return Provider{
		cid: fmt.Sprintf("%v", id),
	}
}

func (p *Provider) GetClientProviders(app, cluster, level, namespace string) (ClientProvider, error) {
	// Default Provider type
	var providerType string = "k8s"

	result := strings.SplitN(cluster, "+", 2)
	if len(result) != 2 {
		log.Error("Invalid cluster name format::", log.Fields{"cluster": cluster})
		return nil, pkgerrors.New("Invalid cluster name format")
	}

	kc, err := utils.GetKubeConfig(cluster, level, namespace)
	if err != nil {
		return nil, err
	}

	if len(kc) > 0 {
		providerType = "k8s"
	} else {
		c, err := utils.GetGitOpsConfig(cluster, level, namespace)
		if err != nil {
			return nil, err
		}
		providerType = c.Props.GitOpsType
		if providerType == "" {
			return nil, pkgerrors.New("No provider type specified")
		}
	}

	switch providerType {
	case "k8s":
		cl, err := k8s.NewK8sProvider(p.cid, app, cluster, level, namespace)
		if err != nil {
			return nil, err
		}
		return cl, nil

	case "fluxcd":
		cl, err := fluxv2.NewFluxv2Provider(p.cid, app, cluster, level, namespace)
		if err != nil {
			return nil, err
		}
		return cl, nil
		//Add other types like Azure Arc, Anthos etc here
	}
	return nil, pkgerrors.New("Provider type not supported")
}
