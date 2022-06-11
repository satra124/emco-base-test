// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
	k8spluginv1alpha1 "gitlab.com/project-emco/core/emco-base/src/monitor/pkg/apis/k8splugin/v1alpha1"
	"log"
	"os"
	"strings"
	"sync"
)

type GithubAccessClient struct {
	cl           gitprovider.Client
	gitUser      string
	gitRepo      string
	cluster      string
	githubDomain string
}

var GitHubClient GithubAccessClient

func SetupGitHubClient() error {
	var err error
	GitHubClient, err = NewGitHubClient()
	return err
}

func NewGitHubClient() (GithubAccessClient, error) {

	githubDomain := "github.com"
	gitUser := os.Getenv("GIT_USERNAME")
	gitToken := os.Getenv("GIT_TOKEN")
	gitRepo := os.Getenv("GIT_REPO")
	clusterName := os.Getenv("GIT_CLUSTERNAME")

	// If any value is not provided then can't store in Git location
	if len(gitRepo) <= 0 || len(gitToken) <= 0 || len(gitUser) <= 0 || len(clusterName) <= 0 {
		log.Printf("Github information not found:: Skipping Github storage")
		return GithubAccessClient{}, nil
	}
	log.Println("GitHub Info found", "gitRepo::", gitRepo, "cluster::", clusterName)

	cl, err := github.NewClient(gitprovider.WithOAuth2Token(gitToken), gitprovider.WithDestructiveAPICalls(true))
	if err != nil {
		return GithubAccessClient{}, err
	}
	return GithubAccessClient{
		cl:           cl,
		gitUser:      gitUser,
		gitRepo:      gitRepo,
		githubDomain: githubDomain,
		cluster:      clusterName,
	}, nil
}

var mutex = sync.Mutex{}

func (c *GithubAccessClient) CommitCRToGitHub(cr *k8spluginv1alpha1.ResourceBundleState, l map[string]string) error {

	// Check if Github Client is available
	if c.cl == nil {
		return nil
	}
	resBytes, err := json.Marshal(cr)
	if err != nil {
		log.Println("json Marshal error for resource::", cr, err)
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

	userRef := gitprovider.UserRef{
		Domain:    c.githubDomain,
		UserLogin: c.gitUser,
	}
	// Create the repo reference
	userRepoRef := gitprovider.UserRepositoryRef{
		UserRef:        userRef,
		RepositoryName: c.gitRepo,
	}
	s := string(resBytes)
	var files []gitprovider.CommitFile
	files = append(files, gitprovider.CommitFile{
		Path:    &path,
		Content: &s,
	})
	commitMessage := "Adding Status for " + path

	// Only one process to commit to Github location to avoid conflicts
	mutex.Lock()
	defer mutex.Unlock()
	userRepo, err := c.cl.UserRepositories().Get(context.Background(), userRepoRef)
	if err != nil {
		return err
	}
	//Commit file to this repo to a branch status
	_, err = userRepo.Commits().Create(context.Background(), "main", commitMessage, files)
	if err != nil {
		return err
	}
	return nil
}
