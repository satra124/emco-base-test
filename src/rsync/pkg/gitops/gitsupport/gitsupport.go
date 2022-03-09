// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package gitsupport

import (
	"context"
	"fmt"
	"strings"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"

	"github.com/fluxcd/go-git-providers/gitprovider"
	pkgerrors "github.com/pkg/errors"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/db"
	emcogit "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogit"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/internal/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type GitProvider struct {
	Cid       string
	Cluster   string
	App       string
	Namespace string
	Level     string
	GitType   string
	GitToken  string
	UserName  string
	Branch    string
	RepoName  string
	Url       string
	Client    gitprovider.Client
}

/*
	Function to create a New gitProivider
	params : cid, app, cluster, level, namespace string
	return : GitProvider, error
*/
func NewGitProvider(cid, app, cluster, level, namespace string) (*GitProvider, error) {

	result := strings.SplitN(cluster, "+", 2)

	c, err := utils.GetGitOpsConfig(cluster, level, namespace)
	if err != nil {
		return nil, err
	}
	// Read from database
	ccc := db.NewCloudConfigClient()
	refObject, err := ccc.GetClusterSyncObjects(result[0], c.Props.GitOpsReferenceObject)

	if err != nil {
		log.Error("Invalid refObject :", log.Fields{"refObj": c.Props.GitOpsReferenceObject, "error": err})
		return nil, err
	}

	kvRef := refObject.Spec.Kv

	var gitType, gitToken, branch, userName, repoName string

	for _, kvpair := range kvRef {
		log.Info("kvpair", log.Fields{"kvpair": kvpair})
		v, ok := kvpair["gitType"]
		if ok {
			gitType = fmt.Sprintf("%v", v)
			continue
		}
		v, ok = kvpair["gitToken"]
		if ok {
			gitToken = fmt.Sprintf("%v", v)
			continue
		}
		v, ok = kvpair["repoName"]
		if ok {
			repoName = fmt.Sprintf("%v", v)
			continue
		}
		v, ok = kvpair["userName"]
		if ok {
			userName = fmt.Sprintf("%v", v)
			continue
		}
		v, ok = kvpair["branch"]
		if ok {
			branch = fmt.Sprintf("%v", v)
			continue
		}
	}
	if len(gitType) <= 0 || len(gitToken) <= 0 || len(branch) <= 0 || len(userName) <= 0 || len(repoName) <= 0 {
		log.Error("Missing information for Git", log.Fields{"gitType": gitType, "token": gitToken, "branch": branch, "userName": userName, "repoName": repoName})
		return nil, pkgerrors.Errorf("Missing Information for Git")
	}

	p := GitProvider{
		Cid:       cid,
		App:       app,
		Cluster:   cluster,
		Level:     level,
		Namespace: namespace,
		GitType:   gitType,
		GitToken:  gitToken,
		Branch:    branch,
		UserName:  userName,
		RepoName:  repoName,
		Url:       "https://" + gitType + ".com/" + userName + "/" + repoName,
	}
	client, err := emcogit.CreateClient(gitToken, gitType)
	if err != nil {
		log.Error("Error getting git client", log.Fields{"err": err})
		return nil, err
	}
	p.Client = client.(gitprovider.Client)
	return &p, nil
}

/*
	Function to get path of files stored in git
	params : string
	return : string
*/

func (p *GitProvider) GetPath(t string) string {
	return "clusters/" + p.Cluster + "/" + t + "/" + p.Cid + "/app/" + p.App + "/"
}

/*
	Function to create a new resource if the not already existing
	params : name string, ref interface{}, content []byte
	return : interface{}, error
*/
func (p *GitProvider) Create(name string, ref interface{}, content []byte) (interface{}, error) {

	path := p.GetPath("context") + name + ".yaml"
	ref = emcogit.Add(path, string(content), ref, p.GitType)
	return ref, nil
}

/*
	Function to apply resource to the cluster
	params : name string, ref interface{}, content []byte
	return : interface{}, error
*/
func (p *GitProvider) Apply(name string, ref interface{}, content []byte) (interface{}, error) {

	//Decode the yaml to create a runtime.Object
	unstruct := &unstructured.Unstructured{}
	//Ignore the returned obj as we expect the data in unstruct
	_, err := utils.DecodeYAMLData(string(content), unstruct)
	if err != nil {
		return nil, err
	}
	// Set Namespace
	unstruct.SetNamespace(p.Namespace)
	b, err := unstruct.MarshalJSON()
	if err != nil {
		return nil, err
	}

	path := p.GetPath("context") + name + ".yaml"
	ref = emcogit.Add(path, string(b), ref, p.GitType)
	return ref, nil

}

/*
	Function to delete resource from the cluster
	params : name string, ref interface{}, content []byte
	return : interface{}, error
*/
func (p *GitProvider) Delete(name string, ref interface{}, content []byte) (interface{}, error) {

	path := p.GetPath("context") + name + ".yaml"
	ref = emcogit.Delete(path, ref, p.GitType)
	return ref, nil

}

/*
	Function to get resource from the cluster
	params : name string, gvkRes []byte
	return : []byte, error
*/
func (p *GitProvider) Get(name string, gvkRes []byte) ([]byte, error) {

	return []byte{}, nil
}

/*
	Function to commit resources to the cluster
	params : ctx context.Context, ref interface{}
	return : error
*/
func (p *GitProvider) Commit(ctx context.Context, ref interface{}) error {

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
	err := emcogit.CommitFiles(ctx, p.Client, p.UserName, p.RepoName, p.Branch, "Commit for "+p.GetPath("context"), ref.([]gitprovider.CommitFile), p.GitType)

	return err
}

/*
	Function for cluster reachablity test
	params : null
	return : error
*/
func (p *GitProvider) IsReachable() error {
	return nil
}
