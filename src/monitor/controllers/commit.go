// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	k8spluginv1alpha1 "gitlab.com/project-emco/core/emco-base/src/monitor/pkg/apis/k8splugin/v1alpha1"
	gitsupport "gitlab.com/project-emco/core/emco-base/src/monitor/pkg/gitops/gitsupport"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"
)

var mutex = sync.Mutex{}

type GithubAccessClient struct {
	gitUser string
	gitRepo string
	cluster string

	gitProvider gitsupport.GitProvider
}

var GitHubClient GithubAccessClient

func SetupGitHubClient() error {
	var err error
	GitHubClient, err = NewGitHubClient()
	return err
}

func NewGitHubClient() (GithubAccessClient, error) {

	gitUser := os.Getenv("GIT_USERNAME")
	gitRepo := os.Getenv("GIT_REPO")
	clusterName := os.Getenv("GIT_CLUSTERNAME")

	// If any value is not provided then can't store in Git location
	if len(gitRepo) <= 0 || len(gitUser) <= 0 || len(clusterName) <= 0 {
		log.Info("Github information not found:: Skipping Github storage", log.Fields{})
		return GithubAccessClient{}, nil
	}
	log.Info("GitHub Info found", log.Fields{"gitRepo::": gitRepo, "cluster::": clusterName})

	gitProvider, err := gitsupport.NewGitProvider()
	if err != nil {
		return GithubAccessClient{}, err
	}

	p := GithubAccessClient{
		gitUser:     gitUser,
		gitRepo:     gitRepo,
		cluster:     clusterName,
		gitProvider: *gitProvider,
	}

	return p, nil
}

func (c *GithubAccessClient) CommitCRToGitHub(cr *k8spluginv1alpha1.ResourceBundleState, l map[string]string) error {

	resBytes, err := json.Marshal(cr)
	if err != nil {
		log.Info("json Marshal error for resource::", log.Fields{"cr": cr, "err": err})
		return err
	}
	// Get cid and app id
	v, ok := l["emco/deployment-id"]
	if !ok {
		return fmt.Errorf("Unexpected error:: Inconsistent labels %v", l)
	}
	result := strings.SplitN(v, "-", 2)
	if len(result) != 2 {
		return fmt.Errorf("Unexpected error:: Inconsistent labels %v", l)
	}
	app := result[1]
	cid := result[0]
	path := "clusters/" + c.cluster + "/status/" + cid + "/app/" + app + "/" + v

	// Add files for commit
	var files interface{}
	files, err = c.gitProvider.Apply(path, files, resBytes)
	if err != nil {
		log.Error("Error in Applying files", log.Fields{"err": err, "files": files, "path": path})
		return err
	}
	branchName := c.cluster

	//commit files
	commitMessage := "Adding Status for " + path + " to branch " + branchName

	// commitfiles
	mutex.Lock()
	defer mutex.Unlock()
	err = c.gitProvider.CommitFilesToBranch(commitMessage, branchName, files)
	if err != nil {
		log.Error("ApplyConfig:: Commit files err", log.Fields{"err": err, "files": files})
		return err
	}

	return nil
}

//function to delete status folder for git
func (c *GithubAccessClient) DeleteStatusFromGit(appName string) error {

	s := strings.SplitN(appName, "-", 2)
	cid := s[0]
	app := s[1]
	path := "clusters/" + c.cluster + "/status/" + cid + "/app/" + app + "/" + appName
	statusBranchName := c.cluster

	var files interface{}
	files, err := c.gitProvider.Delete(path, files)
	if err != nil {
		log.Error("Error in Applying files", log.Fields{"err": err, "files": files, "path": path})
		return err
	}

	//commit files
	commitMessage := "Deleting status for " + appName

	// commitfiles
	mutex.Lock()
	defer mutex.Unlock()
	err = c.gitProvider.CommitFilesToBranch(commitMessage, statusBranchName, files)
	if err != nil {
		log.Error("ApplyConfig:: Commit files err", log.Fields{"err": err, "files": files})
		return err
	}

	return err

}
