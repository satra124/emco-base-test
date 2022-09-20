// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package gitsupport

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"

	git "github.com/libgit2/git2go/v33"
	pkgerrors "github.com/pkg/errors"
	v1alpha1 "gitlab.com/project-emco/core/emco-base/src/monitor/pkg/apis/k8splugin/v1alpha1"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/db"
	emcogit "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogit"
	emcogit2go "gitlab.com/project-emco/core/emco-base/src/rsync/pkg/gitops/emcogit2go"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/internal/utils"
	"gitlab.com/project-emco/core/emco-base/src/rsync/pkg/status"
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
	Client    interface{}
}

/*
	Function to create a New gitProvider
	params : cid, app, cluster, level, namespace string
	return : GitProvider, error
*/
func NewGitProvider(ctx context.Context, cid, app, cluster, level, namespace string) (*GitProvider, error) {

	result := strings.SplitN(cluster, "+", 2)

	c, err := utils.GetGitOpsConfig(ctx, cluster, "0", "default")
	if err != nil {
		return nil, err
	}
	// Read from database
	ccc := db.NewCloudConfigClient()
	refObject, err := ccc.GetClusterSyncObjects(ctx, result[0], c.Props.GitOpsReferenceObject)

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
	client, err := emcogit.CreateClient(userName, gitToken, gitType)
	if err != nil {
		log.Error("Error getting git client", log.Fields{"err": err})
		return nil, err
	}
	p.Client = client
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
	folderName := "/tmp/" + p.UserName + "-" + p.RepoName
	files := emcogit2go.Add(folderName+"/"+path, path, string(content), ref)
	return files, nil
}

/*
	Function to apply resource to the cluster
	params : name string, ref interface{}, content []byte
	return : interface{}, error
*/
func (p *GitProvider) Apply(ctx context.Context, name string, ref interface{}, content []byte) (interface{}, error) {

	path := p.GetPath("context") + name + ".yaml"
	folderName := "/tmp/" + p.UserName + "-" + p.RepoName
	files := emcogit2go.Add(folderName+"/"+path, path, string(content), ref)
	return files, nil

}

/*
	Function to delete resource from the cluster
	params : name string, ref interface{}, content []byte
	return : interface{}, error
*/
func (p *GitProvider) Delete(name string, ref interface{}, content []byte) (interface{}, error) {

	path := p.GetPath("context") + name + ".yaml"
	folderName := "/tmp/" + p.UserName + "-" + p.RepoName
	files := emcogit2go.Delete(folderName+"/"+path, path, ref)
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
func (p *GitProvider) Commit(ctx context.Context, ref interface{}) error {

	var exists bool
	switch ref.(type) {
	case []emcogit2go.CommitFile:
		exists = true
	default:
		exists = false

	}
	// Check for rf
	if !exists {
		log.Error("Commit: No ref found", log.Fields{})
		return nil
	}
	appName := p.Cid + "-" + p.App
	folderName := "/tmp/" + p.UserName + "-" + p.RepoName
	err := emcogit2go.CommitFiles(p.Url, "Commit for "+p.GetPath("context"), appName, folderName, p.UserName, p.GitToken, ref.([]emcogit2go.CommitFile))
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

// Wait time between reading git status (seconds)
var waitTime int = 60
var mutex = sync.Mutex{}

// StartClusterWatcher watches for CR changes in git location
// go routine starts and reads after waitTime
// Thread exists when the AppContext is deleted
func (p *GitProvider) StartClusterWatcher(ctx context.Context) error {
	// obtain the sha key for main
	//obtain shaake
	// ctx := context.Background()

	// latestSHA, err := emcogit.GetLatestCommitSHA(ctx, p.Client, p.UserName, p.RepoName, p.Branch, "", p.GitType)
	// if err != nil {
	// 	return err
	// }
	// // create branch for status
	// branch := p.Cluster + "-" + p.Cid + "-" + p.App
	// err = emcogit.CreateBranch(ctx, p.Client, latestSHA, p.UserName, p.RepoName, branch, p.GitType)
	// if err != nil {
	// 	if !strings.Contains(err.Error(), "422 Reference already exists") {
	// 		return err
	// 	}
	// }

	// foldername for status
	folderName := "/tmp/" + p.Cluster + "-status"
	// Start thread to sync monitor CR
	var lastCommitSHA *git.Oid
	go func() error {
		ctx := context.Background()
		for {
			select {
			case <-time.After(time.Duration(waitTime) * time.Second):
				if ctx.Err() != nil {
					return ctx.Err()
				}
				// Check if AppContext doesn't exist then exit the thread
				if _, err := utils.NewAppContextReference(ctx, p.Cid); err != nil {
					// Delete the Status CR updated by Monitor running on the cluster
					log.Info("Deleting cluster StatusCR", log.Fields{})
					p.DeleteClusterStatusCR(ctx)
					// AppContext deleted - Exit thread
					return nil
				}
				path := p.GetPath("status")
				// branch to track
				branch := p.Cluster

				// // git pull
				mutex.Lock()
				err := emcogit2go.GitPull(p.Url, folderName, branch)
				mutex.Unlock()

				if err != nil {
					log.Error("Error in Pulling Branch", log.Fields{"err": err})
					continue
				}

				//obtain the latest commit SHA
				mutex.Lock()
				latestCommitSHA, err := emcogit2go.GetLatestCommit(folderName, branch)
				mutex.Unlock()
				if err != nil {
					log.Error("Error in obtaining latest commit SHA", log.Fields{"err": err})
				}

				fmt.Println("Latest value")
				fmt.Println(latestCommitSHA)

				fmt.Println("Last Value")
				fmt.Println(lastCommitSHA)

				if lastCommitSHA != latestCommitSHA || lastCommitSHA == nil {
					log.Debug("New Status File, pulling files", log.Fields{"LatestSHA": latestCommitSHA, "LastSHA": lastCommitSHA})
					files, err := emcogit2go.GetFilesInPath(folderName + "/" + path)
					if err != nil {
						log.Debug("Status file not available", log.Fields{"error": err, "resource": path})
						continue
					}
					log.Info("Files to track status", log.Fields{"files": files})

					if len(files) > 0 {
						// Only one file expected in the location
						fileContent, err := emcogit2go.GetFileContent(files[0])
						if err != nil {
							log.Error("", log.Fields{"error": err, "cluster": p.Cluster, "resource": path})
							return err
						}
						content := &v1alpha1.ResourceBundleState{}
						_, err = utils.DecodeYAMLData(fileContent, content)
						if err != nil {
							log.Error("", log.Fields{"error": err, "cluster": p.Cluster, "resource": path})
							return err
						}
						status.HandleResourcesStatus(ctx, p.Cid, p.App, p.Cluster, content)
					}

					lastCommitSHA = latestCommitSHA

				}

			// Check if the context is canceled
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}()
	return nil
}

// DeleteClusterStatusCR deletes the status CR provided by the monitor on the cluster
func (p *GitProvider) DeleteClusterStatusCR(ctx context.Context) error {
	// Delete the status CR
	// branch to track
	// branch := p.Cluster + "-" + p.Cid + "-" + p.App

	// ctx := context.Background()

	// // Delete the branch
	// err := emcogit.DeleteBranch(ctx, p.Client, p.UserName, p.RepoName, branch, p.GitType)
	// if err != nil {
	// 	log.Error("Error in deleting branch", log.Fields{"err": err})
	// 	return err
	// }

	//delete the dummy branch as well
	// folderName := "/tmp/" + p.Cluster + "-" + p.Cid
	// // // // open a repo
	// repo, err := git.OpenRepository(folderName)
	// if err != nil {
	// 	log.Error("Error in Opening the git repository", log.Fields{"err": err})
	// 	return err
	// }
	// err = emcogit2go.DeleteBranch(repo, branch)
	// if err != nil {
	// 	return err
	// }

	// err = emcogit2go.PushDeleteBranch(repo, branch)
	// if err != nil {
	// 	return err
	// }
	// remove the local folder
	// err := os.RemoveAll(folderName)
	// if err != nil {
	// 	return err
	// }
	return nil
}
