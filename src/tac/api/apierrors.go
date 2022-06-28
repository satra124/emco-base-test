// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package api

import (
	"net/http"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
)

var apiErrors = []apierror.APIError{
	{ID: "workflow Hook Intent Already exists", Message: "workflow hook intent already exists", Status: http.StatusConflict},
	{ID: "empty Post Body", Message: "empty post body", Status: http.StatusBadRequest},
	{ID: "error decoding json body", Message: "error decoding json body", Status: http.StatusUnprocessableEntity},
	{ID: "Workflow Hook not found", Message: "Workflow Hook not found", Status: http.StatusNotFound},
}

// HandleErrors exposes the generic action controller API errors
func HandleErrors(params map[string]string, err error, mod interface{}) apierror.APIError {
	return apierror.HandleErrors(params, err, mod, apiErrors)
}
