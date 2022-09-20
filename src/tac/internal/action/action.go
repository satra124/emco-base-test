// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package action

import (
	"github.com/pkg/errors"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/appcontext"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/grpc/contextupdate"
	"gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	con "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/utils"
	"gitlab.com/project-emco/core/emco-base/src/tac/pkg/module"
	eta "gitlab.com/project-emco/core/emco-base/src/workflowmgr/pkg/emcotemporalapi"

	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// UpdateAppContext applies the supplied intent against the given AppContext ID
func UpdateAppContext(intentName, appContextId, hookType string) error {
	var hookString string
	if len(hookType) == 0 {
		hookString = "pre-install"
	} else {
		hookString = "pre-update"
	}

	var ac appcontext.AppContext

	_, err := ac.LoadAppContext(context.Background(), appContextId)
	if err != nil {
		logutils.Error("Failed to get the appContext.",
			logutils.Fields{
				"ID": appContextId})
		return errors.Wrapf(err, "Failed to get the appContext with ID: %s.", appContextId)
	}

	_, err = ac.GetCompositeAppHandle(context.Background())
	if err != nil {
		return err
	}

	appContext, err := con.ReadAppContext(context.Background(), appContextId)
	if err != nil {
		logutils.Error("Failed to get the compositeApp for the appContext.",
			logutils.Fields{
				"ID": appContextId})
		return errors.Wrapf(err, "Failed to get the compositeApp for the appContext with ID: %s", appContextId)
	}

	project := appContext.CompMetadata.Project
	app := appContext.CompMetadata.CompositeApp
	version := appContext.CompMetadata.Version
	group := appContext.CompMetadata.DeploymentIntentGroup

	// look up workflow pre install hooks
	pre, err := module.NewClient().WorkflowIntentClient.GetSpecificHooks(project, app, version, group, hookString)
	if err != nil {
		logutils.Error("Failed to get the intents for the deploymentIntentGroup.",
			logutils.Fields{
				"DeploymentIntentGroup": group})
		return errors.Wrapf(err, "Failed to get the intents for the deploymentIntentGroup: %s", group)
	}

	// iterate through all pre-install hooks
	logutils.Info("UpdateAppContext TAC ... Iterating through all workflow preinstall intents for this deployment intent group.",
		logutils.Fields{})

	for _, p := range pre {

		// log which pre-install hook we are working on
		logutils.Info("UpdateAppContext TAC ... Starting pre-install hooks.",
			logutils.Fields{
				"Temporal hook": p.Metadata.Name,
				"Project":       project,
				"App":           app,
				"Version":       version,
				"DIG":           group,
			})

		// execute the workflow
		err = RunWorkflow(p.Spec.WfClientSpec.WfClientEndpointName,
			strconv.Itoa(p.Spec.WfClientSpec.WfClientEndpointPort), p.Spec.WfTemporalSpec.WfClientName, p.Spec.WfTemporalSpec)
		if err != nil {
			logutils.Error("Failed to run workflow in pre-install hooks.",
				logutils.Fields{
					"Hook name": p.Metadata.Name})
			return errors.Wrapf(err, "Failed to run workflow in pre-install hook: %s", p.Metadata.Name)
		}

	}

	// log that pre-install is finished
	logutils.Info("UpdateAppContext TAC ... Pre-install hooks finished.",
		logutils.Fields{
			"Project": project,
			"App":     app,
			"Version": version,
			"DIG":     group,
		})

	return nil
}

// TerminateAppContext is the to run the pre/post terminate workflows for this deployment intent group
func TerminateAppContext(appContextId string) error {
	var ac appcontext.AppContext

	logutils.Info("Terminate TAC ... pre-terminate hooks starting.", logutils.Fields{})

	_, err := ac.LoadAppContext(context.Background(), appContextId)
	if err != nil {
		logutils.Error("Failed to get the appContext.",
			logutils.Fields{
				"ID": appContextId})
		return errors.Wrapf(err, "Failed to get the appContext with ID: %s.", appContextId)
	}

	_, err = ac.GetCompositeAppHandle(context.Background())
	if err != nil {
		return err
	}

	appContext, err := con.ReadAppContext(context.Background(), appContextId)
	if err != nil {
		logutils.Error("Failed to get the compositeApp for the appContext.",
			logutils.Fields{
				"ID": appContextId})
		return errors.Wrapf(err, "Failed to get the compositeApp for the appContext with ID: %s", appContextId)
	}

	project := appContext.CompMetadata.Project
	app := appContext.CompMetadata.CompositeApp
	version := appContext.CompMetadata.Version
	group := appContext.CompMetadata.DeploymentIntentGroup

	// look up workflow pre install hooks
	pre, err := module.NewClient().WorkflowIntentClient.GetSpecificHooks(project, app, version, group, "pre-termination")
	if err != nil {
		logutils.Error("Failed to get the intents for the deploymentIntentGroup.",
			logutils.Fields{
				"DeploymentIntentGroup": group})
		return errors.Wrapf(err, "Failed to get the intents for the deploymentIntentGroup: %s", group)
	}

	// iterate through all pre-install hooks
	logutils.Info("Terminate TAC ... Iterating through all workflow pre-terminate intents for this deployment intent group.",
		logutils.Fields{})

	for _, p := range pre {

		// log which pre-install hook we are working on
		logutils.Info("Terminate TAC ... Starting pre-terminate hooks.",
			logutils.Fields{
				"Temporal hook": p.Metadata.Name,
				"Project":       project,
				"App":           app,
				"Version":       version,
				"DIG":           group,
			})

		// execute the workflow
		err = RunWorkflow(p.Spec.WfClientSpec.WfClientEndpointName,
			strconv.Itoa(p.Spec.WfClientSpec.WfClientEndpointPort), p.Spec.WfTemporalSpec.WfClientName, p.Spec.WfTemporalSpec)
		if err != nil {
			logutils.Error("Failed to run workflow in pre-install hooks.",
				logutils.Fields{
					"Hook name": p.Metadata.Name})
			return errors.Wrapf(err, "Failed to run workflow in pre-install hook: %s", p.Metadata.Name)
		}

	}

	// log that pre-install is finished
	logutils.Info("Terminate TAC ... Pre-terminate hooks finished.",
		logutils.Fields{
			"Project": project,
			"App":     app,
			"Version": version,
			"DIG":     group,
		})

	return nil
}

func PostEvent(appContextId string, et contextupdate.EventType) error {
	eventString := []string{"post-install", "post-termination", "post-update"}
	var ac appcontext.AppContext

	logutils.Info("PostEvent TAC ... "+eventString[et]+" hooks starting.", logutils.Fields{})

	_, err := ac.LoadAppContext(context.Background(), appContextId)
	if err != nil {
		logutils.Error("PostEvent Failed to get the appContext.",
			logutils.Fields{
				"ID": appContextId})
		return errors.Wrapf(err, "PostEvent Failed to get the appContext with ID: %s.", appContextId)
	}

	_, err = ac.GetCompositeAppHandle(context.Background())
	if err != nil {
		return err
	}

	appContext, err := con.ReadAppContext(context.Background(), appContextId)
	if err != nil {
		logutils.Error("PostEvent Failed to get the compositeApp for the appContext.",
			logutils.Fields{
				"ID": appContextId})
		return errors.Wrapf(err, "PostEvent Failed to get the compositeApp for the appContext with ID: %s", appContextId)
	}

	project := appContext.CompMetadata.Project
	app := appContext.CompMetadata.CompositeApp
	version := appContext.CompMetadata.Version
	group := appContext.CompMetadata.DeploymentIntentGroup

	// look up workflow pre install hooks
	post, err := module.NewClient().WorkflowIntentClient.GetSpecificHooks(project, app, version, group, eventString[et])
	if err != nil {
		logutils.Error("PostEvent Failed to get the intents for the deploymentIntentGroup.",
			logutils.Fields{
				"DeploymentIntentGroup": group})
		return errors.Wrapf(err, "PostEvent Failed to get the intents for the deploymentIntentGroup: %s", group)
	}

	// iterate through all pre-install hooks
	logutils.Info("PostEvent TAC ... Iterating through all workflow "+eventString[et]+" intents for this deployment intent group.",
		logutils.Fields{})

	for _, p := range post {

		// log which pre-install hook we are working on
		logutils.Info("PostEvent TAC ... Starting "+eventString[et]+" hooks.",
			logutils.Fields{
				"Temporal hook": p.Metadata.Name,
				"Project":       project,
				"App":           app,
				"Version":       version,
				"DIG":           group,
			})

		// execute the workflow
		err = RunWorkflow(p.Spec.WfClientSpec.WfClientEndpointName,
			strconv.Itoa(p.Spec.WfClientSpec.WfClientEndpointPort), p.Spec.WfTemporalSpec.WfClientName, p.Spec.WfTemporalSpec)
		if err != nil {
			logutils.Error("Failed to run workflow in "+eventString[et]+" hooks.",
				logutils.Fields{
					"Hook name": p.Metadata.Name})
			return errors.Wrapf(err, "Failed to run workflow in PostEvent hook: %s", p.Metadata.Name)
		}

	}

	// log that PostEvent is finished
	logutils.Info("PostEvent TAC ... "+eventString[et]+" hooks finished.",
		logutils.Fields{
			"Project": project,
			"App":     app,
			"Version": version,
			"DIG":     group,
		})

	return nil
}

func RunWorkflow(clientEndpoint, clientPort, clientName string, spec eta.WfTemporalSpec) error {
	url := "http://" + clientEndpoint + ":" +
		clientPort + "/invoke/" +
		clientName

	jsonBytes, err := json.Marshal(spec)
	if err != nil {
		logutils.Error("Temporal Action Controller.... Error marshalling the workflow temporal spec into post body",
			logutils.Fields{"error": err.Error()})
		return err
	}
	logutils.Info("Temporal Action Controller.... Run workflow",
		logutils.Fields{"url": url, "workflow temporal spec": string(jsonBytes)})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		logutils.Error("Temporal Action Controller.... Run workflow could not post",
			logutils.Fields{"error": err.Error()})
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		logutils.Error("Temporal Action Controller.... Run workflow POST returned error",
			logutils.Fields{"status code": resp.Status, "urL": url})
		return errors.New("Error starting workflow")
	}

	return nil
}
