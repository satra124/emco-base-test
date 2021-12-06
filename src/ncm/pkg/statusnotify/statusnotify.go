// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

package statusnotify

import (
	pkgerrors "github.com/pkg/errors"
	clustermodule "gitlab.com/project-emco/core/emco-base/src/clm/pkg/cluster"
	"gitlab.com/project-emco/core/emco-base/src/ncm/pkg/scheduler"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/appcontext"
	statusnotifypb "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc/statusnotify"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc/statusnotifyserver"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/state"
)

type clusterHelpers struct{}

func getClusterKeyValues(reg *statusnotifypb.StatusRegistration) (string, string, error) {
	key := reg.GetKey()

	var clusterKey *statusnotifypb.ClusterKey

	switch key.(type) {
	case *statusnotifypb.StatusRegistration_ClusterKey:
		clusterKey = reg.GetClusterKey()
	default:
		return "", "", pkgerrors.New("Status Notification Registration - Key is not a Cluster key")
	}

	return clusterKey.GetClusterProvider(), clusterKey.GetCluster(), nil
}

func (d clusterHelpers) GetAppContextId(reg *statusnotifypb.StatusRegistration) (string, error) {
	clusterProvider, cluster, err := getClusterKeyValues(reg)
	if err != nil {
		return "", err
	}

	si, err := clustermodule.NewClusterClient().GetClusterState(clusterProvider, cluster)
	if err != nil {
		return "", pkgerrors.Wrapf(err, "Cluster state not found: %v+%v", clusterProvider, cluster)
	}

	return state.GetStatusContextIdFromStateInfo(si), nil
}

func (d clusterHelpers) PrepareStatusNotification(reg *statusnotifypb.StatusRegistration) *statusnotifypb.StatusNotification {
	n := new(statusnotifypb.StatusNotification)

	clusterProvider, cluster, err := getClusterKeyValues(reg)
	if err != nil {
		return n
	}

	statusType, output, apps, clusters, resources := statusnotifyserver.GetStatusParameters(reg)

	statusResult, err := scheduler.NewSchedulerClient().NetworkIntentsStatus(clusterProvider, cluster, "", statusType, output, apps, clusters, resources)
	if err != nil {
		return n
	}

	if statusResult.Status == appcontext.AppContextStatusEnum.Instantiated {
		switch reg.StatusType {
		case statusnotifypb.StatusValue_DEPLOYED:
			n.StatusValue = statusnotifypb.StatusValue_DEPLOYED
		case statusnotifypb.StatusValue_READY:
			// TODO:  calling the cluster READY in this case is an assumption.  'monitor' does not currently
			// return information about cluster resources, so the real status is unknown.
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

	// NOTE: skip details for type READY, as 'monitor' does not currently return cluster network resource information.
	if reg.Output == statusnotifypb.OutputType_ALL && reg.StatusType == statusnotifypb.StatusValue_DEPLOYED {
		details := make([]*statusnotifypb.StatusDetail, 0)

		clusterStatus := statusnotifypb.ClusterStatus{}
		clusterStatus.Cluster = statusResult.Cluster.Cluster
		clusterStatus.ClusterProvider = statusResult.Cluster.ClusterProvider
		resources := make([]*statusnotifypb.ResourceStatus, 0)
		clusterDeployed := true
		for _, resource := range statusResult.Cluster.Resources {
			resourceStatus := statusnotifypb.ResourceStatus{}
			resourceStatus.Name = resource.Name
			resourceStatus.Gvk = &statusnotifypb.GVK{Group: resource.Gvk.Group, Version: resource.Gvk.Version, Kind: resource.Gvk.Kind}
			if resource.RsyncStatus == "Applied" {
				resourceStatus.StatusValue = statusnotifypb.StatusValue_DEPLOYED
			} else {
				resourceStatus.StatusValue = statusnotifypb.StatusValue_NOT_DEPLOYED
				clusterDeployed = false
			}
			resources = append(resources, &resourceStatus)
		}
		if clusterDeployed {
			clusterStatus.StatusValue = statusnotifypb.StatusValue_DEPLOYED
		} else {
			clusterStatus.StatusValue = statusnotifypb.StatusValue_NOT_DEPLOYED
		}
		clusterStatus.Resources = resources
		details = append(details, &statusnotifypb.StatusDetail{StatusDetail: &statusnotifypb.StatusDetail_Cluster{Cluster: &clusterStatus}})
		n.Details = details
	}

	return n
}

func StartStatusNotifyServer() *statusnotifyserver.StatusNotifyServer {
	return statusnotifyserver.NewStatusNotifyServer("clusterStatus", clusterHelpers{})
}
