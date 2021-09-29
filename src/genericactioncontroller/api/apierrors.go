package api

import (
	"net/http"

	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/apierror"
)

var apiErrors = []apierror.APIError{
	{ID: "Customization already exists", Message: "Customization already exists", Status: http.StatusConflict},
	{ID: "Customization not found", Message: "Customization not found", Status: http.StatusNotFound},
	{ID: "Customization Spec File Content not found", Message: "Customization Spec File Content not found", Status: http.StatusNotFound},
	{ID: "GenericK8sIntent already exists", Message: "GenericK8sIntent already exists", Status: http.StatusConflict},
	{ID: "Generic K8s Intent not found", Message: "Generic K8s Intent not found", Status: http.StatusNotFound},
	{ID: "Resource already exists", Message: "Resource already exists", Status: http.StatusConflict},
	{ID: "Resource not found", Message: "Resource not found", Status: http.StatusNotFound},
	{ID: "Resource File Content not found", Message: "Resource File Content not found", Status: http.StatusNotFound},
}
