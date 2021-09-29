// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package api

import (
	"github.com/gorilla/mux"
	moduleLib "gitlab.com/project-emco/core/emco-base/src/genericactioncontroller/pkg/module"

	//"log"
	"fmt"
	"reflect"
)

var moduleClient *moduleLib.Client

// For the given client and testClient, if the testClient is not null and
// implements the client manager interface corresponding to client, then
// return the testClient, otherwise return the client.
func setClient(client, testClient interface{}) interface{} {
	switch cl := client.(type) {
	case *moduleLib.GenericK8sIntentClient:
		if testClient != nil && reflect.TypeOf(testClient).Implements(reflect.TypeOf((*moduleLib.GenericK8sIntentManager)(nil)).Elem()) {
			c, ok := testClient.(moduleLib.GenericK8sIntentManager)
			if ok {
				return c
			}
		}
	case *moduleLib.ResourceClient:
		if testClient != nil && reflect.TypeOf(testClient).Implements(reflect.TypeOf((*moduleLib.ResourceManager)(nil)).Elem()) {
			c, ok := testClient.(moduleLib.ResourceManager)
			if ok {
				return c
			}
		}

	case *moduleLib.CustomizationClient:
		if testClient != nil && reflect.TypeOf(testClient).Implements(reflect.TypeOf((*moduleLib.CustomizationManager)(nil)).Elem()) {
			c, ok := testClient.(moduleLib.ResourceManager)
			if ok {
				return c
			}
		}
	default:
		fmt.Printf("unknown type %T\n", cl)
	}
	return client
}

// NewRouter returns the mux router after plugging in all the handlers
func NewRouter(testClient interface{}) *mux.Router {
	moduleClient = moduleLib.NewClient()

	router := mux.NewRouter().PathPrefix("/v2").Subrouter()

	genericK8sintentHandler := generick8sintentHandler{
		client: setClient(moduleClient.GenericK8sIntent, testClient).(moduleLib.GenericK8sIntentManager),
	}

	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents", genericK8sintentHandler.createHandler).Methods("POST")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents", genericK8sintentHandler.getHandler).Methods("GET")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}", genericK8sintentHandler.getHandler).Methods("GET")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}", genericK8sintentHandler.putHandler).Methods("PUT")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}", genericK8sintentHandler.deleteHandler).Methods("DELETE")

	baseResourceHandler := resourceHandler{
		client: setClient(moduleClient.BaseResource, testClient).(moduleLib.ResourceManager),
	}

	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources", baseResourceHandler.createResourceHandler).Methods("POST")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources", baseResourceHandler.getResourceHandler).Methods("GET")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources/{genericResource}", baseResourceHandler.getResourceHandler).Methods("GET")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources/{genericResource}", baseResourceHandler.putResourceHandler).Methods("PUT")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources/{genericResource}", baseResourceHandler.deleteResourceHandler).Methods("DELETE")

	customizationHandler := customizationHandler{
		client: setClient(moduleClient.Customization, testClient).(moduleLib.CustomizationManager),
	}

	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources/{genericResource}/customizations", customizationHandler.createCustomizationHandler).Methods("POST")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources/{genericResource}/customizations", customizationHandler.getCustomizationHandler).Methods("GET")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources/{genericResource}/customizations/{customization}", customizationHandler.getCustomizationHandler).Methods("GET")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources/{genericResource}/customizations/{customization}", customizationHandler.putCustomizationHandler).Methods("PUT")
	router.HandleFunc("/projects/{project}/composite-apps/{compositeApp}/{compositeAppVersion}/deployment-intent-groups/{deploymentIntentGroup}/generic-k8s-intents/{genericK8sIntent}/resources/{genericResource}/customizations/{customization}", customizationHandler.deleteCustomizationHandler).Methods("DELETE")

	return router
}
