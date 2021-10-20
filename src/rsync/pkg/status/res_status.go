// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package status

import (
	"encoding/json"
	rb "gitlab.com/project-emco/core/emco-base/src/monitor/pkg/apis/k8splugin/v1alpha1"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/appcontext"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/depend"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/grpc/readynotifyserver"
)

// Update status for the App ready on a cluster and check if app ready on all clusters
func HandleResourcesStatus(acID, app, cluster string, rbData *rb.ResourceBundleState) {

	// Look up the contextId
	var ac appcontext.AppContext
	_, err := ac.LoadAppContext(acID)
	if err != nil {
		log.Info("::App context not found::", log.Fields{"acID": acID, "app": app, "cluster": cluster, "err": err})
		return
	}

	// Produce yaml representation of the status
	vjson, err := json.Marshal(rbData.Status)
	if err != nil {
		log.Info("::Error marshalling status information::", log.Fields{"acID": acID, "app": app, "cluster": cluster, "err": err})
		return
	}

	chandle, err := ac.GetClusterHandle(app, cluster)
	if err != nil {
		log.Info("::Error getting cluster handle::", log.Fields{"acID": acID, "app": app, "cluster": cluster, "err": err})
		return
	}
	// Get the handle for the context/app/cluster status object
	handle, _ := ac.GetLevelHandle(chandle, "status")

	// If status handle was not found, then create the status object in the appcontext
	if handle == nil {
		ac.AddLevelValue(chandle, "status", string(vjson))
	} else {
		ac.UpdateStatusValue(handle, string(vjson))
	}

	if CheckAppReadyStatus(acID, app, cluster, rbData) {
		// Inform Rsync dependency management
		go depend.ResourcesReady(acID, app, cluster)
	}

	// Send notification to the subscribers
	err = readynotifyserver.SendAppContextNotification(acID)
	if err != nil {
		log.Error("::Error sending ReadyNotify to subscribers::", log.Fields{"acID": acID, "app": app, "cluster": cluster, "err": err})
	}
}

// Check if all the resources are ready on a cluster
func IsAppReady(rbData *rb.ResourceBundleState) bool {

	readyChecker := NewReadyChecker(PausedAsReady(true), CheckJobs(true))
	var avail bool = false
	for _, s := range rbData.Status.ServiceStatuses {
		avail = true
		if !readyChecker.ServiceReady(&s) {
			return false
		}
	}

	for _, d := range rbData.Status.DeploymentStatuses {
		avail = true
		if !readyChecker.DeploymentReady(&d) {
			return false
		}
	}

	for _, d := range rbData.Status.DaemonSetStatuses {
		avail = true
		if !readyChecker.DaemonSetReady(&d) {
			return false
		}
	}

	for _, j := range rbData.Status.JobStatuses {
		avail = true
		if !readyChecker.JobReady(&j) {
			return false
		}
	}

	for _, s := range rbData.Status.StatefulSetStatuses {
		avail = true
		if !readyChecker.StatefulSetReady(&s) {
			return false
		}
	}
	for _, p := range rbData.Status.PodStatuses {
		avail = true
		if !readyChecker.PodReady(&p) {
			return false
		}
	}
	if !avail {
		log.Info("No resources found in monitor CR", log.Fields{"rbData": rbData})
		return false
	}

	return true

}

func CheckAppReadyStatus(acID, app string, cluster string, rbData *rb.ResourceBundleState) bool {

	ac := appcontext.AppContext{}
	_, err := ac.LoadAppContext(acID)
	if err != nil {
		return false
	}
	// If Application is ready on the cluster, Update AppContext
	// If the application is not ready stop processing
	if !IsAppReady(rbData) {
		setClusterResourcesReady(ac, app, cluster, false)
		return false
	}
	log.Info("ClusterResourcesReady App is ready on cluster", log.Fields{"acID": acID, "app": app, "cluster": cluster})

	setClusterResourcesReady(ac, app, cluster, true)

	// Check if all the clusters are ready
	cl, err := ac.GetClusterNames(app)
	if err != nil {
		return false
	}
	for _, cn := range cl {
		if !getClusterResourcesReady(ac, app, cn) {
			// Some cluster is not ready
			return false
		}
	}
	return true
}

// setClusterResourceReady sets the cluster ready status
func setClusterResourcesReady(ac appcontext.AppContext, app, cluster string, value bool) error {

	ch, err := ac.GetClusterHandle(app, cluster)
	if err != nil {
		return err
	}
	rsh, _ := ac.GetLevelHandle(ch, "resourcesready")
	// If resource ready handle was not found, then create it
	if rsh == nil {
		ac.AddLevelValue(ch, "resourcesready", value)
	} else {
		ac.UpdateStatusValue(rsh, value)
	}
	return nil
}

// getClusterResourceReady gets the cluster ready status
func getClusterResourcesReady(ac appcontext.AppContext, app, cluster string) bool {
	ch, err := ac.GetClusterHandle(app, cluster)
	if err != nil {
		return false
	}
	rsh, _ := ac.GetLevelHandle(ch, "resourcesready")
	if rsh != nil {
		status, err := ac.GetValue(rsh)
		if err != nil {
			return false
		}
		return status.(bool)
	}
	return false
}
