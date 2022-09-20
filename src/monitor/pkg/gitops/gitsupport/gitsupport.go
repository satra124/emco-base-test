// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package gitsupport

import (
	"os"
	"strings"

	pkgerrors "github.com/pkg/errors"
	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"

	//emcogit "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogit"
	emcogit2go "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogit2go"
	emcogithub "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogithub"
)

type GitProvider struct {
	// Cid       string
	// Cluster   string
	// App       string
	// Namespace string
	// Level     string
	// GitType   string
	// GitToken  string
	// UserName  string
	// Branch    string
	// RepoName  string
	// Url       string

	GitUser      string
	GitRepo      string
	GitToken     string
	Cluster      string
	GitType      string
	GitBranch    string
	Url          string
	gitInterface GitInterfaceProvider
}

type GitInterfaceProvider interface {
	AddToCommit(fileName, folderName, content string, ref interface{}) interface{}
	DeleteToCommit(fileName, folderName string, ref interface{}) interface{}
	CommitFiles(app, message, folderName string, files interface{}) error
	ClusterWatcher(cid, app, cluster string, waitTime int) error
	CommitFilesToBranch(commitMessage, branchName, folderName string, files interface{}) error
}

/*
	Function to create a New gitProvider
	params : cid, app, cluster, level, namespace string
	return : GitProvider, error
*/
func NewGitProvider() (*GitProvider, error) {

	// result := strings.SplitN(cluster, "+", 2)

	// c, err := utils.GetGitOpsConfig(cluster, "0", "default")
	// if err != nil {
	// 	return nil, err
	// }
	// // Read from database
	// ccc := db.NewCloudConfigClient()
	// refObject, err := ccc.GetClusterSyncObjects(result[0], c.Props.GitOpsReferenceObject)

	// if err != nil {
	// 	log.Error("Invalid refObject :", log.Fields{"refObj": c.Props.GitOpsReferenceObject, "error": err})
	// 	return nil, err
	// }

	// kvRef := refObject.Spec.Kv

	// var gitType, gitToken, branch, userName, repoName, url string

	// for _, kvpair := range kvRef {
	// 	log.Info("kvpair", log.Fields{"kvpair": kvpair})
	// 	v, ok := kvpair["gitType"]
	// 	if ok {
	// 		gitType = fmt.Sprintf("%v", v)
	// 		continue
	// 	}
	// 	v, ok = kvpair["gitToken"]
	// 	if ok {
	// 		gitToken = fmt.Sprintf("%v", v)
	// 		continue
	// 	}
	// 	v, ok = kvpair["repoName"]
	// 	if ok {
	// 		repoName = fmt.Sprintf("%v", v)
	// 		continue
	// 	}
	// 	v, ok = kvpair["userName"]
	// 	if ok {
	// 		userName = fmt.Sprintf("%v", v)
	// 		continue
	// 	}
	// 	v, ok = kvpair["branch"]
	// 	if ok {
	// 		branch = fmt.Sprintf("%v", v)
	// 		continue
	// 	}
	// 	v, ok = kvpair["url"]
	// 	if ok {
	// 		url = fmt.Sprintf("%v", v)
	// 		continue
	// 	}
	// }
	// if len(gitType) <= 0 || len(gitToken) <= 0 || len(branch) <= 0 || len(userName) <= 0 || len(repoName) <= 0 || len(url) <= 0 {
	// 	log.Error("Missing information for Git", log.Fields{"gitType": gitType, "token": gitToken, "branch": branch, "userName": userName, "repoName": repoName, "url": url})
	// 	return nil, pkgerrors.Errorf("Missing Information for Git")
	// }

	gitBranch := os.Getenv("GIT_BRANCH")
	gitType := os.Getenv("GIT_TYPE")
	gitUser := os.Getenv("GIT_USERNAME")
	gitToken := os.Getenv("GIT_TOKEN")
	gitRepo := os.Getenv("GIT_REPO")
	clusterName := os.Getenv("GIT_CLUSTERNAME")
	gitUrl := os.Getenv("GIT_URL")

	// If any value is not provided then can't store in Git location
	if len(gitRepo) <= 0 || len(gitToken) <= 0 || len(gitUser) <= 0 || len(clusterName) <= 0 || len(gitUrl) <= 0 {
		log.Info("Github information not found:: Skipping Github storage", log.Fields{})
		return nil, pkgerrors.Errorf("Missing Information for Git")
	}
	log.Info("GitHub Info found", log.Fields{"gitRepo::": gitRepo, "cluster::": clusterName})

	p := GitProvider{
		GitUser:   gitUser,
		GitRepo:   gitRepo,
		GitToken:  gitToken,
		GitType:   gitType,
		GitBranch: gitBranch,
		Cluster:   clusterName,
		Url:       gitUrl,
	}

	if strings.EqualFold(gitType, "github") {
		p.gitInterface, _ = emcogithub.NewGithub(p.Cluster, p.Url, p.GitBranch, p.GitUser, p.GitRepo, p.GitToken)
	} else {
		p.gitInterface = emcogit2go.NewGit2Go(p.Url, p.GitBranch, p.GitUser, p.GitRepo, p.GitToken)
	}

	return &p, nil
}

/*
	Function to get path of files stored in git
	params : string
	return : string
*/

// func (p *GitProvider) GetPath(t string) string {
// 	return "clusters/" + p.Cluster + "/" + t + "/" + p.Cid + "/app/" + p.App + "/"
// }

// /*
// 	Function to create a new resource if the not already existing
// 	params : name string, ref interface{}, content []byte
// 	return : interface{}, error
// */
// func (p *GitProvider) Create(name string, ref interface{}, content []byte) (interface{}, error) {

// 	path := p.GetPath("context") + name + ".yaml"
// 	files := p.gitInterface.AddToCommit(path, string(content), ref)
// 	return files, nil
// }

/*
	Function to apply resource to the cluster
	params : name string, ref interface{}, content []byte
	return : interface{}, error
*/
func (p *GitProvider) Apply(path string, ref interface{}, content []byte) (interface{}, error) {

	folderName := "/tmp/" + p.GitUser + "-" + p.GitRepo + "-" + p.Cluster
	files := p.gitInterface.AddToCommit(path, folderName, string(content), ref)
	return files, nil

}

/*
	Function to delete resource from the cluster
	params : name string, ref interface{}, content []byte
	return : interface{}, error
*/
func (p *GitProvider) Delete(path string, ref interface{}, content []byte) (interface{}, error) {

	folderName := "/tmp/" + p.GitUser + "-" + p.GitRepo + "-" + p.Cluster
	files := p.gitInterface.DeleteToCommit(path, folderName, ref)
	return files, nil

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
func (p *GitProvider) CommitFiles(commitMessage, branchName, folderName string, files interface{}) error {

	err := p.gitInterface.CommitFilesToBranch(commitMessage, branchName, folderName, files)
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
