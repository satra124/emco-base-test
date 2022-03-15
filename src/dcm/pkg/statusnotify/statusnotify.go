// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package statusnotify

import (
	pkgerrors "github.com/pkg/errors"
	dcm "gitlab.com/project-emco/core/emco-base/src/dcm/pkg/module"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/appcontext"
	statusnotifypb "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc/statusnotify"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc/statusnotifyserver"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/state"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/status"
)

type lcHelpers struct{}

func getLcKeyValues(reg *statusnotifypb.StatusRegistration) (string, string, error) {
	key := reg.GetKey()

	var lcKey *statusnotifypb.LcKey

	switch key.(type) {
	case *statusnotifypb.StatusRegistration_LcKey:
		lcKey = reg.GetLcKey()
	default:
		return "", "", pkgerrors.New("Status Notification Registration - Key is not a Deployment Intent Group key")
	}

	return lcKey.GetProject(), lcKey.GetLogicalCloud(), nil
}

func (d lcHelpers) GetAppContextId(reg *statusnotifypb.StatusRegistration) (string, error) {
	p, lc, err := getLcKeyValues(reg)
	if err != nil {
		return "", err
	}

	si, err := dcm.NewLogicalCloudClient().GetState(p, lc)
	if err != nil {
		return "", pkgerrors.Wrap(err, "Logical cloud state not found: "+lc)
	}

	return state.GetStatusContextIdFromStateInfo(si), nil
}

func (d lcHelpers) StatusQuery(reg *statusnotifypb.StatusRegistration, qStatusInstance, qType, qOutput string, fApps, fClusters, fResources []string) status.StatusResult {
	p, lc, err := getLcKeyValues(reg)
	if err != nil {
		return status.StatusResult{}
	}

	statusResult, err := dcm.NewLogicalCloudClient().GenericStatus(p, lc, qStatusInstance, qType, qOutput, fClusters, fResources)
	if err != nil {
		return status.StatusResult{}
	}
	return statusResult
}

func (d lcHelpers) PrepareStatusNotification(reg *statusnotifypb.StatusRegistration, statusResult status.StatusResult) *statusnotifypb.StatusNotification {
	n := new(statusnotifypb.StatusNotification)

	// TODO: use when logical cloud status supports these filter parameters
	//statusType, output, apps, clusters, resources := statusnotifyserver.GetStatusParameters(reg)
	// TODO: fix up once dcm more fully supports the status query

	if statusResult.Status == appcontext.AppContextStatusEnum.Instantiated {
		switch reg.StatusType {
		case statusnotifypb.StatusValue_DEPLOYED:
			n.StatusValue = statusnotifypb.StatusValue_DEPLOYED
		case statusnotifypb.StatusValue_READY:
			n.StatusValue = statusnotifypb.StatusValue_READY
		}
	} else {
		switch reg.StatusType {
		case statusnotifypb.StatusValue_DEPLOYED:
			n.StatusValue = statusnotifypb.StatusValue_NOT_DEPLOYED
		case statusnotifypb.StatusValue_READY:
			n.StatusValue = statusnotifypb.StatusValue_NOT_READY
		}
	}

	return n
}

func StartStatusNotifyServer() *statusnotifyserver.StatusNotifyServer {
	return statusnotifyserver.NewStatusNotifyServer("digStatus", lcHelpers{})
}
