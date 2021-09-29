// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/project-emco/core/emco-base/src/dcm/pkg/module"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
)

// userPermissionHandler is used to store backend implementations objects
type userPermissionHandler struct {
	client module.UserPermissionManager
}

// CreateHandler handles creation of the user permission entry in the database
func (h userPermissionHandler) createHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	project := vars["project"]
	logicalCloud := vars["logicalCloud"]
	var v module.UserPermission

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

	// User-Permission Name is required.
	if v.MetaData.UserPermissionName == "" {
		msg := "Missing name in POST request"
		log.Error(msg, log.Fields{})
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	ret, err := h.client.CreateUserPerm(project, logicalCloud, v)
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

// getAllHandler handles GET operations over user permissions
// Returns a list of User Permissions
func (h userPermissionHandler) getAllHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	logicalCloud := vars["logicalCloud"]
	var ret interface{}
	var err error

	ret, err = h.client.GetAllUserPerms(project, logicalCloud)
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
// Returns a User Permission
func (h userPermissionHandler) getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	logicalCloud := vars["logicalCloud"]
	name := vars["userPermission"]
	var ret interface{}
	var err error

	ret, err = h.client.GetUserPerm(project, logicalCloud, name)
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

// UpdateHandler handles Update operations on a particular user permission
func (h userPermissionHandler) updateHandler(w http.ResponseWriter, r *http.Request) {
	var v module.UserPermission
	vars := mux.Vars(r)
	project := vars["project"]
	logicalCloud := vars["logicalCloud"]
	name := vars["userPermission"]

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

	// User-Permission Name is required.
	if v.MetaData.UserPermissionName == "" {
		log.Error("API: Missing name in PUT request", log.Fields{})
		http.Error(w, "Missing name in PUT request", http.StatusBadRequest)
		return
	}

	ret, err := h.client.UpdateUserPerm(project, logicalCloud, name, v)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, v, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		log.Error(err.Error(), log.Fields{})
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
}

//deleteHandler handles DELETE operations on a particular record
func (h userPermissionHandler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	logicalCloud := vars["logicalCloud"]
	name := vars["userPermission"]

	err := h.client.DeleteUserPerm(project, logicalCloud, name)
	if err != nil {
		apiErr := apierror.HandleErrors(vars, err, nil, apiErrors)
		http.Error(w, apiErr.Message, apiErr.Status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
