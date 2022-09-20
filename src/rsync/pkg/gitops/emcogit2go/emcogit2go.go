package emcogit2go

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"

	git "github.com/libgit2/git2go/v33"
)

const (
	maxrand = 0x7fffffffffffffff
)

var mutex = sync.Mutex{}

func getCredCallBack(userName, token string) func(url string, username string, allowedTypes git.CredType) (*git.Credential, error) {
	return func(url string, username string, allowedTypes git.CredType) (*git.Credential, error) {
		username = userName
		password := token
		cred, err := git.NewCredUserpassPlaintext(username, password)
		return cred, err
	}
}

// CommitFile contains high-level information about a file added to a commit.
type CommitFile struct {
	// Path is path where this file is located.
	// +required
	Add bool `json:"add"`

	Path *string `json:"path"`

	FileName *string `json:"filename"`

	// Content is the content of the file.
	// +required
	Content *string `json:"content,omitempty"`
}

//function to delete a file
func deleteFile(filenName string) error {
	// Removing file from the directory
	// Using Remove() function
	err := os.Remove(filenName)
	if err != nil {
		log.Error("Error in Deleting file from the tmp folder", log.Fields{"err": err})
		return err
	}

	return nil
}

//function to create a new file
func createFile(fileName string, content string) error {
	if err := os.MkdirAll(filepath.Dir(fileName), 0770); err != nil {
		return err
	}

	f, err := os.Create(fileName)

	if err != nil {
		log.Error("Error in Creating file in the tmp folder", log.Fields{"err": err})
		return err
	}

	defer f.Close()

	_, err2 := f.WriteString(content)

	if err2 != nil {
		log.Error("Error in writing file from the tmp folder", log.Fields{"err": err2})
		return err2
	}
	return nil
}

// function to commit files to a branch
func CommitFiles(url, message, appName, folderName, userName, token string, files []CommitFile) error {

	mutex.Lock()
	defer mutex.Unlock()

	var repo *git.Repository
	// // // open a repo
	//clone the repo
	// clone git the repo to local repo
	check, err := Exists(folderName)

	if !check {
		if err := os.Mkdir(folderName, os.ModePerm); err != nil {
			log.Error("Error in creating the dir", log.Fields{"Error": err})
			return err
		}
		// // clone the repo
		repo, err = git.Clone(url, folderName, &git.CloneOptions{CheckoutBranch: "main", CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutSafe}})
		if err != nil {
			log.Error("Error cloning the repo", log.Fields{"Error": err})
			return err
		}

	}
	repo, err = git.OpenRepository(folderName)
	if err != nil {
		log.Error("Error in Opening the git repository", log.Fields{"err": err, "appName": appName})
		return err
	}

	signature := &git.Signature{
		Name:  "Adarsh Vincent",
		Email: "a.v@gmail.com",
		When:  time.Now(),
	}
	// //create a new branch from main
	// ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	// rn := ra.Int63n(maxrand)
	// id := fmt.Sprintf("%v", rn)

	// branchName := appName + "-" + id

	// create the new branch (May cause problems, try to get the headCommit of main)
	//checkout the new branch
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + "main")
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": "main"})
		return err
	}

	branch, err := repo.References.Lookup("refs/heads/" + "main")
	if err != nil {
		log.Info("Error in looking up ref", log.Fields{"err": err})
		return err
	}

	// head, err := repo.Head()
	// if err != nil {
	// 	log.Error("Error in obtaining the head of the repo", log.Fields{"err": err})
	// 	return err
	// }

	// headCommit, err := repo.LookupCommit(head.Target())
	// if err != nil {
	// 	log.Error("Error in obtainging the head commit", log.Fields{"err": err, "headCommit": headCommit})
	// 	return err
	// }
	// branch, err := repo.CreateBranch(branchName, headCommit, false)
	// if err != nil {
	// 	log.Error("Error in Creating branch", log.Fields{"err": err, "branchName": branchName, "headCommit": headCommit, "branch": branch})
	// 	return err
	// }

	// //checkout the new branch
	// // set head to point to the created branch
	// err = repo.SetHead("refs/heads/" + branchName)
	// if err != nil {
	// 	log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
	// 	return err
	// }

	//Update the index with files and obtain the latest index
	//loop through all files and update the index
	idx, err := repo.Index()
	if err != nil {
		log.Error("Error in obtaining the repo index", log.Fields{"err": err, "idx": idx})
		return err
	}

	for _, file := range files {
		if file.Add {
			idx, err = addToCommit(idx, *file.Path, *file.FileName, *file.Content)
		} else {
			idx, err = deleteFromCommit(idx, *file.Path, *file.FileName)
		}

		if err != nil {
			log.Error("Error in adding or deleting file to commit", log.Fields{"err": err, "idx": idx})
			return err
		}
	}
	//commit the files to the branch
	treeId, err := idx.WriteTree()
	if err != nil {
		log.Error("Error from idx.WriteTree()", log.Fields{"err": err})
		return err
	}

	err = idx.Write()
	if err != nil {
		log.Error("Error in Deleting file from idx.Write()", log.Fields{"err": err})
		return err
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		log.Error("Error in looking up tree", log.Fields{"err": err, "treeId": treeId})
		return err
	}

	commitTarget, err := repo.LookupCommit(branch.Target())
	if err != nil {
		log.Error("Error in Looking up Commit for commit", log.Fields{"err": err})
		return err
	}

	branchName := "main"
	_, err = repo.CreateCommit("refs/heads/"+branchName, signature, signature, message, tree, commitTarget)
	if err != nil {
		log.Error("Error in creating a commit", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	// //merge branch to main
	// err = mergeToMain(repo, branchName, signature)
	// //set head to master
	// err = repo.SetHead("refs/heads/main")
	// //delete the created branch
	// err = DeleteBranch(repo, branchName)
	//push master to origin remote
	err = PushBranch(repo, "main", userName, token)

	if err != nil {
		return err
	}

	return nil
}

// function to push branch to remote origin
func PushBranch(repo *git.Repository, branchName, userName, token string) error {
	// push the branch to origin
	remote, err := repo.Remotes.Lookup("origin")

	cbs := &git.RemoteCallbacks{
		CredentialsCallback: getCredCallBack(userName, token),
	}

	err = remote.Push([]string{"+refs/heads/" + branchName}, &git.PushOptions{RemoteCallbacks: *cbs})
	if err != nil {
		log.Error("Error in Pushing the branch", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	return nil
}

// function to push branch to remote origin
// func PushDeleteBranch(repo *git.Repository, branchName string) error {
// 	// push the branch to origin
// 	remote, err := repo.Remotes.Lookup("origin")

// 	cbs := &git.RemoteCallbacks{
// 		CredentialsCallback: getCredCallBack(userName, token),
// 	}

// 	err = remote.Push([]string{":refs/heads/" + branchName}, &git.PushOptions{RemoteCallbacks: *cbs})
// 	if err != nil {
// 		log.Error("Error in Pushing the branch", log.Fields{"err": err, "branchName": branchName})
// 		return err
// 	}

// 	return nil
// }

//function to merge branch to main (Should include a commit as well)
func mergeToMain(repo *git.Repository, branchName string, signature *git.Signature) error {
	// get reference for the target merge branch
	mergeBranch, err := repo.References.Lookup("refs/heads/" + branchName)
	if err != nil {
		log.Error("Error in obtaining the reference for branch to merge to main", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	mergeHeadMergeBranch, err := repo.AnnotatedCommitFromRef(mergeBranch)
	if err != nil {
		log.Error("Error in obtaining the head of the branch to merge", log.Fields{"err": err, "mergeHeadMergeBranch": mergeHeadMergeBranch})
		return err
	}
	mergeHeads := make([]*git.AnnotatedCommit, 1)

	mergeHeads[0] = mergeHeadMergeBranch

	err = repo.Merge(mergeHeads, nil, nil)
	if err != nil {
		log.Error("Error in Merging the branch", log.Fields{"err": err, "mergeHeads": mergeHeads})
		return err
	}

	mergeMessage, err := repo.Message()
	if err != nil {
		return err
	}

	log.Debug("Merge Message", log.Fields{"mergeMessage": mergeMessage})

	err = commitMergeToMaster(repo, signature, "Merge commit to main")
	if err != nil {
		log.Error("Error in commit Merge to main", log.Fields{"err": err})
		return err
	}

	return nil
}

//function to delete branch
func DeleteBranch(repo *git.Repository, branchName string) error {
	branchA, err := repo.LookupBranch(branchName, git.BranchLocal)
	err = branchA.Delete()
	if err != nil {
		return err
	}
	return nil
}

//function to check if folder exists
// exists returns whether the given file or directory exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// function to add file for commit
func addToCommit(idx *git.Index, path, fileName, contents string) (*git.Index, error) {

	err := createFile(path, contents)
	if err != nil {
		return nil, err
	}
	// add file to staging area
	err = idx.AddByPath(fileName)
	if err != nil {
		return nil, err
	}
	return idx, nil

}

// function to delete file for commit
func deleteFromCommit(idx *git.Index, path, fileName string) (*git.Index, error) {

	// check if the file Exists
	check, err := Exists(path)
	if err != nil {
		return nil, err
	}
	if check {
		err := idx.RemoveByPath(fileName)
		if err != nil {
			return nil, err
		}

		err = deleteFile(path)
		if err != nil {
			return nil, err
		}
	}
	return idx, nil
}

// function to add file to commit files array
func Add(path, fileName, content string, ref interface{}) []CommitFile {
	files := append(convertToCommitFile(ref), CommitFile{
		Add:      true,
		Path:     &path,
		FileName: &fileName,
		Content:  &content,
	})

	return files

}

// function to delete file from commit files array
func Delete(path, fileName string, ref interface{}) []CommitFile {
	files := append(convertToCommitFile(ref), CommitFile{
		Add:      false,
		Path:     &path,
		FileName: &fileName,
	})

	return files
}

func commitMergeToMaster(repo *git.Repository, signature *git.Signature, message string) error {
	//commit the merge to main
	branchName := "main"
	idx, err := repo.Index()
	if err != nil {
		log.Error("commitMergeToMaster: Error in obtaining the repo index", log.Fields{"err": err, "idx": idx})
		return err
	}

	branchMain, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return err
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		log.Error("commitMergeToMaster: Error from idx.WriteTree()", log.Fields{"err": err})
		return err
	}

	err = idx.Write()
	if err != nil {
		log.Error("commitMergeToMaster: Error in Deleting file from idx.Write()", log.Fields{"err": err})
		return err
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		log.Error("commitMergeToMaster: Error in looking up tree", log.Fields{"err": err, "treeId": treeId})
		return err
	}

	commitTarget, err := repo.LookupCommit(branchMain.Target())
	if err != nil {
		log.Error("commitMergeToMaster: Error in Looking up Commit for commit", log.Fields{"err": err})
		return err
	}

	_, err = repo.CreateCommit("refs/heads/"+branchName, signature, signature, message, tree, commitTarget)
	if err != nil {
		log.Error("commitMergeToMaster:Error in creating a commit", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	return nil
}

/*
	Helper function to convert interface to []gitprovider.CommitFile
	params: files interface{}
	return: []gitprovider.CommitFile
*/
func convertToCommitFile(ref interface{}) []CommitFile {
	var exists bool
	switch ref.(type) {
	case []CommitFile:
		exists = true
	default:
		exists = false
	}
	var rf []CommitFile
	// Create rf is doesn't exist
	if !exists {
		rf = []CommitFile{}
	} else {
		rf = ref.([]CommitFile)
	}
	return rf
}

//function to create branch
func CreateBranch(folderName, branchName string) (*git.Branch, error) {
	// // // open a repo
	repo, err := git.OpenRepository(folderName)
	if err != nil {
		log.Error("Error in Opening the git repository", log.Fields{"err": err})
		return nil, err
	}

	// create the new branch (May cause problems, try to get the headCommit of main)
	//checkout the new branch
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + "main")
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		log.Error("Error in obtaining the head of the repo", log.Fields{"err": err})
		return nil, err
	}

	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		log.Error("Error in obtainging the head commit", log.Fields{"err": err, "headCommit": headCommit})
		return nil, err
	}
	branch, err := repo.CreateBranch(branchName, headCommit, false)
	if err != nil {
		log.Error("Error in Creating branch", log.Fields{"err": err, "branchName": branchName, "headCommit": headCommit, "branch": branch})
		return nil, err
	}

	return branch, nil
}

//function to pull branch
func GitPull(url, folderName, branchName string) error {
	check, err := Exists(folderName)

	if !check {
		if err := os.Mkdir(folderName, os.ModePerm); err != nil {
			return err
		}
		// // clone the repo
		_, err := git.Clone(url, folderName, &git.CloneOptions{CheckoutBranch: branchName, CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutSafe}})
		if err != nil {
			log.Error("Error in cloning repo Git Pull", log.Fields{"err": err})
			// remove the local folder
			err := os.RemoveAll(folderName)
			return err
		}
	}
	// // // open a repo
	repo, err := git.OpenRepository(folderName)
	if err != nil {
		log.Error("Error in Opening the git repository", log.Fields{"err": err})
		return err
	}

	// Locate remote
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		log.Error("Error in looking up Origin", log.Fields{"err": err})
		return err
	}

	// Fetch changes from remote
	if err := remote.Fetch([]string{}, nil, ""); err != nil {
		log.Error("Error in Fetching", log.Fields{"err": err})
		return err
	}

	// check if a local branch exist if not do a checkout
	localBranch, err := repo.LookupBranch(branchName, git.BranchLocal)
	// No local branch, lets checkout
	if localBranch == nil || err != nil {
		err = CheckoutBranch(folderName, branchName)
		if err != nil {
			log.Error("Git checkout error", log.Fields{"err": err, "branchName": branchName})
			return err
		}
		return nil
	}
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + branchName)
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	// Get remote master
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + branchName)
	if err != nil {
		return err
	}

	remoteBranchID := remoteBranch.Target()
	// Get annotated commit
	annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)

	if err != nil {
		return err
	}

	// Do the merge analysis
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = annotatedCommit
	analysis, _, err := repo.MergeAnalysis(mergeHeads)
	if err != nil {
		return err
	}

	log.Debug("Git Pull Information ", log.Fields{"remoteBranch": remoteBranch, "remoteBranchID": remoteBranchID, "annotatedCommit": annotatedCommit, "analysis": analysis, "mergeHeads": mergeHeads})

	head, err := repo.Head()
	if err != nil {
		return err
	}

	if analysis&git.MergeAnalysisUpToDate != 0 {
		log.Info("MergeAnalysisUpToDate", log.Fields{"analysis": analysis})
		return nil
	} else if analysis&git.MergeAnalysisNormal != 0 {
		log.Info("MergeAnalysisNormal", log.Fields{"analysis": analysis})

		// Just merge changes
		if err := repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
			log.Error("Merge Error", log.Fields{"err": err})
			return err
		}
		// Check for conflicts
		index, err := repo.Index()
		if err != nil {
			return err
		}

		if index.HasConflicts() {
			return errors.New("Conflicts encountered. Please resolve them.")
		}

		// Make the merge commit
		sig := &git.Signature{Name: "Adarsh", Email: "adarsh@cool", When: time.Now()}
		if err != nil {
			return err
		}

		// Get Write Tree
		treeId, err := index.WriteTree()
		if err != nil {
			return err
		}

		// Get repo head

		tree, err := repo.LookupTree(treeId)
		if err != nil {
			return err
		}

		head, err := repo.Head()
		if err != nil {
			return err
		}
		localCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}

		remoteCommit, err := repo.LookupCommit(remoteBranchID)
		if err != nil {
			return err
		}

		_, err = repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)
		if err != nil {
			return err
		}
		// Clean up
		repo.StateCleanup()
	} else if analysis&git.MergeAnalysisFastForward != 0 {
		log.Info("MergeAnalysisFastForward", log.Fields{"analysis": analysis})
		// Fast-forward changes
		// Get remote tree
		remoteTree, err := repo.LookupTree(remoteBranchID)
		if err != nil {
			return err
		}

		// Checkout
		if err := repo.CheckoutTree(remoteTree, nil); err != nil {
			return err
		}

		branchRef, err := repo.References.Lookup("refs/heads/" + branchName)
		if err != nil {
			return err
		}

		// Point branch to the object
		branchRef.SetTarget(remoteBranchID, "")
		if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
			return err
		}

	} else {
		log.Info("Unexpected merge analysis result", log.Fields{"analysis": analysis})
		return fmt.Errorf("Unexpected merge analysis result %d", analysis)
	}
	return nil
}

//function to checkout a branch
func CheckoutBranch(folderName, branchName string) error {
	repo, err := git.OpenRepository(folderName)
	if err != nil {
		return err
	}
	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}
	//Getting the reference for the remote branch
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + branchName)
	if err != nil {
		log.Error("Failed to find remote branch: ", log.Fields{"branchName": branchName})
		return err
	}
	defer remoteBranch.Free()

	// Lookup for commit from remote branch
	commit, err := repo.LookupCommit(remoteBranch.Target())
	if err != nil {
		log.Error("Failed to find remote branch commit: ", log.Fields{"branchName": branchName})
		return err
	}
	defer commit.Free()

	localBranch, err := repo.LookupBranch(branchName, git.BranchLocal)
	// No local branch, lets create one
	if localBranch == nil || err != nil {
		// Creating local branch
		localBranch, err = repo.CreateBranch(branchName, commit, false)
		if err != nil {
			log.Error("Failed to create local branch: ", log.Fields{"branchName": branchName})
			return err
		}

		// Setting upstream to origin branch
		err = localBranch.SetUpstream("origin/" + branchName)
		if err != nil {
			log.Error("Failed to create upstream to origin/ ", log.Fields{"branchName": branchName})
			return err
		}
	}
	if localBranch == nil {
		return errors.New("Error while locating/creating local branch")
	}
	defer localBranch.Free()

	// Getting the tree for the branch
	localCommit, err := repo.LookupCommit(localBranch.Target())
	if err != nil {
		log.Error("Failed to lookup for commit in local branch ", log.Fields{"branchName": branchName})
		return err
	}
	defer localCommit.Free()

	tree, err := repo.LookupTree(localCommit.TreeId())
	if err != nil {
		log.Error("Failed to lookup for tree", log.Fields{"branchName": branchName})
		return err
	}
	defer tree.Free()

	// Checkout the tree
	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		log.Error("Failed to checkout tree ", log.Fields{"branchName": branchName})
		return err
	}
	return nil
}

//Function to get files in a path
func GetFilesInPath(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		file, err := os.Open(path)
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		//Add only if the path is not a directory
		if !fileInfo.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

//Function to get File Contents
func GetFileContent(filePath string) (string, error) {

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Convert []byte to string and print to screen
	text := string(content)

	return text, nil
}

//function to get the latest commit
func GetLatestCommit(path, branchName string) (*git.Oid, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		log.Error("Error in Opening repo", log.Fields{"err": err, "path": path})
		return nil, err
	}
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + branchName)
	if err != nil {
		log.Error("Error in setting head", log.Fields{"err": err, "branchName": branchName})
		return nil, err
	}
	head, err := repo.Head()
	if err != nil {
		log.Error("Error obtaining head", log.Fields{"err": err, "head": head})
		return nil, err
	}
	headCommit := head.Target()

	return headCommit, nil
}

// git pull using command
func GitPullCMD(url, folderName, branchName string) error {
	mutex.Lock()
	defer mutex.Unlock()
	// // // open a repo
	check, err := Exists(folderName)

	if !check {
		if err := os.Mkdir(folderName, os.ModePerm); err != nil {
			return err
		}
		// // clone the repo
		_, err := git.Clone(url, folderName, &git.CloneOptions{CheckoutBranch: branchName, CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutSafe}})
		if err != nil {
			return err
		}
	}
	repo, err := git.OpenRepository(folderName)
	if err != nil {
		log.Error("Error in Opening the git repository", log.Fields{"err": err})
		return err
	}
	// Locate remote
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		log.Error("Error in looking up Origin", log.Fields{"err": err})
		return err
	}

	// Fetch changes from remote
	if err := remote.Fetch([]string{}, nil, ""); err != nil {
		log.Error("Error in Fetching", log.Fields{"err": err})
		return err
	}
	// // check if a local branch exist if not do a checkout
	localBranch, err := repo.LookupBranch(branchName, git.BranchLocal)
	// // No local branch, lets create one
	if localBranch == nil || err != nil {
		err = CheckoutBranchCMD(folderName, branchName)
		if err != nil {
			log.Error("Error in Fetching", log.Fields{"err": err})
		}
	}
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + branchName)
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
		return err
	}
	cmd := exec.Command("git", "pull")
	cmd.Dir = folderName
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

//function to checkout branch using command
func CheckoutBranchCMD(folderName, branchName string) error {
	_, err := git.OpenRepository(folderName)
	if err != nil {
		return err
	}

	// using git command
	cmd := exec.Command("git", "checkout", "-b", branchName, "origin/"+branchName)
	cmd.Dir = folderName
	err = cmd.Run()
	if err != nil {
		log.Error("Git checkout error", log.Fields{"err": err, "branchName": branchName})
		return err
	}
	return nil
}
