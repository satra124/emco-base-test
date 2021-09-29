// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	dcm "gitlab.com/project-emco/core/emco-base/src/dcm/pkg/module"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	orch "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/module"
)

// logicalCloudHandler is used to store backend implementations objects
type logicalCloudHandler struct {
	client               dcm.LogicalCloudManager
	clusterClient        dcm.ClusterManager
	quotaClient          dcm.QuotaManager
	userPermissionClient dcm.UserPermissionManager
}

// CreateHandler handles the creation of a logical cloud
func (h logicalCloudHandler) createHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	project := vars["project"]
	var v dcm.LogicalCloud

	err := json.NewDecoder(r.Body).Decode(&v)
	switch {
	case err == io.EOF:
		log.Error(err.Error(), log.Fields{})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		log.Error(err.Error(), log.Fields{})
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Logical Cloud Name is required.
	if v.MetaData.LogicalCloudName == "" {
		msg := "Missing name in POST request"
		log.Error(msg, log.Fields{})
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Validate that the specified Project exists
	// before associating a Logical Cloud with it
	p := orch.NewProjectClient()
	_, err = p.GetProject(project)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, v, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	ret, err := h.client.Create(project, v)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, v, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Error(err.Error(), log.Fields{})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getAllHandler handles GET operations over logical clouds
// Returns a list of Logical Clouds
func (h logicalCloudHandler) getAllHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	var ret interface{}
	var err error

	ret, err = h.client.GetAll(project)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Error(err.Error(), log.Fields{})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getHandler handles GET operations on a particular name
// Returns a Logical Cloud
func (h logicalCloudHandler) getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	name := vars["logicalCloud"]
	var ret interface{}
	var err error

	ret, err = h.client.Get(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Error(err.Error(), log.Fields{})
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// updateHandler handles Update operations on a particular logical cloud
func (h logicalCloudHandler) updateHandler(w http.ResponseWriter, r *http.Request) {
	var v dcm.LogicalCloud
	vars := mux.Vars(r)
	project := vars["project"]
	name := vars["logicalCloud"]

	err := json.NewDecoder(r.Body).Decode(&v)
	switch {
	case err == io.EOF:
		log.Error(err.Error(), log.Fields{})
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	case err != nil:
		log.Error(err.Error(), log.Fields{})
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if v.MetaData.LogicalCloudName == "" {
		log.Error("API: Missing name in PUT request", log.Fields{})
		http.Error(w, "Missing name in PUT request", http.StatusBadRequest)
		return
	}

	ret, err := h.client.Update(project, name, v)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, v, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
}

// deleteHandler handles Delete operations on a particular logical cloud
func (h logicalCloudHandler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	name := vars["logicalCloud"]

	err := h.client.Delete(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// instantiateHandler handles instantiateing a particular logical cloud
func (h logicalCloudHandler) instantiateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	name := vars["logicalCloud"]

	// Get logical cloud
	lc, err := h.client.Get(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	// Get Clusters
	clusters, err := h.clusterClient.GetAllClusters(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	// Get Quotas
	quotas, err := h.quotaClient.GetAllQuotas(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	userPermissions, err := h.userPermissionClient.GetAllUserPerms(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	// Instantiate the Logical Cloud
	err = dcm.Instantiate(project, lc, clusters, quotas, userPermissions)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// terminateHandler handles terminating a particular logical cloud
func (h logicalCloudHandler) terminateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	name := vars["logicalCloud"]

	// Get logical cloud
	lc, err := h.client.Get(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	// Get Clusters
	clusters, err := h.clusterClient.GetAllClusters(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	// Get Quotas
	quotas, err := h.quotaClient.GetAllQuotas(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	// Terminate the Logical Cloud
	err = dcm.Terminate(project, lc, clusters, quotas)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// stopHandler handles aborting the pending instantiation or termination of a logical cloud
func (h logicalCloudHandler) stopHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	name := vars["logicalCloud"]

	// Get logical cloud
	lc, err := h.client.Get(project, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	// Stop the instantiation
	err = dcm.Stop(project, lc)
	if err != nil {
		log.Error(err.Error(), log.Fields{})
		if err.Error() == "Logical Clouds can't be stopped" {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
