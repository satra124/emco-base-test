// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package fluxv2

import (
	"context"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/status"

	"github.com/fluxcd/go-git-providers/gitprovider"
	emcogithub "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogithub"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/internal/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (p *Fluxv2Provider) getPath() string {
	return "clusters/" + p.cluster + "/context/" + p.cid + "/app/" + p.app + "/"
}

func convertToCommitFile(ref interface{}) []gitprovider.CommitFile {
	var exists bool
	switch ref.(type) {
	case []gitprovider.CommitFile:
		exists = true
	default:
		exists = false
	}
	var rf []gitprovider.CommitFile
	// Create rf is doesn't exist
	if !exists {
		rf = []gitprovider.CommitFile{}
	} else {
		rf = ref.([]gitprovider.CommitFile)
	}
	return rf
}

// Creates a new resource if the not already existing
func (p *Fluxv2Provider) Create(name string, ref interface{}, content []byte) (interface{}, error) {
	// Add the label based on the Status Appcontext ID
	label := p.cid + "-" + p.app
	b, err := status.TagResource(content, label)
	if err != nil {
		log.Error("Error Tag Resource with label:", log.Fields{"err": err, "label": label, "resource": name})
		return nil, err
	}
	path := p.getPath() + name + ".yaml"
	rf := convertToCommitFile(ref)
	ref = emcogithub.Add(path, string(b), rf)
	return ref, nil
}

// Apply resource to the cluster
func (p *Fluxv2Provider) Apply(name string, ref interface{}, content []byte) (interface{}, error) {
	// Add the label based on the Status Appcontext ID
	label := p.cid + "-" + p.app

	//Decode the yaml to create a runtime.Object
	unstruct := &unstructured.Unstructured{}
	//Ignore the returned obj as we expect the data in unstruct
	_, err := utils.DecodeYAMLData(string(content), unstruct)
	if err != nil {
		return nil, err
	}
	//Add the tracking label to all resources created here
	labels := unstruct.GetLabels()
	//Check if labels exist for this object
	if labels == nil {
		labels = map[string]string{}
	}
	// Set label
	labels["emco/deployment-id"] = label
	unstruct.SetLabels(labels)
	// Set Namespace
	unstruct.SetNamespace(p.namespace)

	// This checks if the resource we are creating has a podSpec in it
	// Eg: Deployment, StatefulSet, Job etc..
	// If a PodSpec is found, the label will be added to it too.
	//connector.TagPodsIfPresent(unstruct, client.GetInstanceID())
	status.TagPodsIfPresent(unstruct, label)
	b, err := unstruct.MarshalJSON()
	if err != nil {
		return nil, err
	}

	path := p.getPath() + name + ".yaml"
	rf := convertToCommitFile(ref)
	ref = emcogithub.Add(path, string(b), rf)
	return ref, nil

}

// Delete resource from the cluster
func (p *Fluxv2Provider) Delete(name string, ref interface{}, content []byte) (interface{}, error) {
	path := p.getPath() + name + ".yaml"
	rf := convertToCommitFile(ref)
	ref = emcogithub.Delete(path, rf)
	return ref, nil

}

// Get resource from the cluster
func (p *Fluxv2Provider) Get(name string, gvkRes []byte) ([]byte, error) {

	return []byte{}, nil
}

// Commit resources to the cluster
func (p *Fluxv2Provider) Commit(ctx context.Context, ref interface{}) error {
	var exists bool
	switch ref.(type) {
	case []gitprovider.CommitFile:
		exists = true
	default:
		exists = false

	}
	// Check for rf
	if !exists {
		log.Error("Commit: No ref found", log.Fields{})
		return nil
	}
	err := emcogithub.CommitFiles(ctx, p.client, p.userName, p.repoName, p.branch, "Commit for "+p.getPath(), ref.([]gitprovider.CommitFile))

	return err
}

// IsReachable cluster reachablity test
func (p *Fluxv2Provider) IsReachable() error {
	return nil
}

func (m *Fluxv2Provider) TagResource(res []byte, label string) ([]byte, error) {
	b, err := status.TagResource(res, label)
	if err != nil {
		log.Error("Error Tag Resoruce with label:", log.Fields{"err": err, "label": label, "resource": res})
		return nil, err
	}
	return b, nil
}
