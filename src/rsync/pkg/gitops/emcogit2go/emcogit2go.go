package emcogit2go

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
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

func credentialsCallback(url string, username string, allowedTypes git.CredType) (*git.Credential, error) {
	username = "chitti-intel"
	password := "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"
	cred, err := git.NewCredUserpassPlaintext(username, password)
	return cred, err
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

// func main() {

// 	// username := "chitti-intel"

// 	// // create a temp dir
// 	// outDir, _ := ioutil.TempDir("", "test-git-")

// 	folderName := "/tmp/newDir"

// 	check, err := Exists(folderName)

// 	if !check {
// 		if err := os.Mkdir(folderName, os.ModePerm); err != nil {
// 			log.Fatal(err)
// 		}
// 		// // clone the repo
// 		repo, err := git.Clone("https://github.com/chitti-intel/arc-k8s-demo", folderName, &git.CloneOptions{CheckoutBranch: "master", CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutSafe}})
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(repo)
// 	}

// 	// fmt.Println(outDir)

// 	// // // open a repo
// 	repo, err := git.OpenRepository(folderName)

// 	// head, err := repo.Head()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// headCommit, err := repo.LookupCommit(head.Target())
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// // create a new branch
// 	// branchName := "git2go-tutorial-v8"
// 	// branch, err := repo.CreateBranch(branchName, headCommit, false)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// // set head to point to the created branch

// 	// err = repo.SetHead("refs/heads/" + branchName)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// idx, err := repo.Index()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	files := []CommitFile{}

// 	// Delete a file
// 	fileName := "testing"
// 	path := "/tmp/newDir/" + fileName

// 	files = Delete(path, fileName, files)

// 	// add another file
// 	//create a new file in the path
// 	fileName = "test-file-new.txt"
// 	path = "/tmp/newDir/" + fileName
// 	contents := "Hello there\n"

// 	files = Add(path, fileName, contents, files)

// 	// commit file to the new branch
// 	CommitFiles(repo, "Add test-file-new and delete testing file", "collectd-app-", files)

// 	// err = idx.AddByPath("test-file-v3.txt")

// 	// //merge branch to main
// 	// mergeToMain(repo, branchName, signature)

// 	// err = repo.SetHead("refs/heads/master")

// 	// // delete branch
// 	// deleteBranch(repo, branchName)

// 	// // // push branch to remote orgin
// 	// pushBranch(repo, "master")

// }

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

	fmt.Println("done")
	return nil
}

// function to commit files to a branch
func CommitFiles(message, appName, folderName string, files []CommitFile) error {

	mutex.Lock()
	defer mutex.Unlock()
	// // // open a repo
	//clone the repo
	if err := os.Mkdir(folderName, os.ModePerm); err != nil {
		return err
	}
	// // clone the repo
	_, err := git.Clone("https://github.com/chitti-intel/test-flux-v3", folderName, &git.CloneOptions{CheckoutBranch: "main", CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutSafe}})
	if err != nil {
		return err
	}
	repo, err := git.OpenRepository(folderName)
	if err != nil {
		log.Error("Error in Opening the git repository", log.Fields{"err": err, "appName": appName})
		return err
	}

	signature := &git.Signature{
		Name:  "Adarsh Vincent",
		Email: "a.v@gmail.com",
		When:  time.Now(),
	}
	//create a new branch from main
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	rn := ra.Int63n(maxrand)
	id := fmt.Sprintf("%v", rn)

	branchName := appName + "-" + id

	// create the new branch (May cause problems, try to get the headCommit of main)
	//checkout the new branch
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + "main")
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	head, err := repo.Head()
	if err != nil {
		log.Error("Error in obtaining the head of the repo", log.Fields{"err": err})
		return err
	}

	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		log.Error("Error in obtainging the head commit", log.Fields{"err": err, "headCommit": headCommit})
		return err
	}
	branch, err := repo.CreateBranch(branchName, headCommit, false)
	if err != nil {
		log.Error("Error in Creating branch", log.Fields{"err": err, "branchName": branchName, "headCommit": headCommit, "branch": branch})
		return err
	}

	//checkout the new branch
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + branchName)
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
		return err
	}

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

	_, err = repo.CreateCommit("refs/heads/"+branchName, signature, signature, message, tree, commitTarget)
	if err != nil {
		log.Error("Error in creating a commit", log.Fields{"err": err, "branchName": branchName})
		return err
	}
	//merge with main
	//merge branch to main
	err = mergeToMain(repo, branchName, signature)
	//set head to master
	err = repo.SetHead("refs/heads/main")
	//delete the created branch
	err = DeleteBranch(repo, branchName)
	//push master to origin remote
	err = PushBranch(repo, "main")

	// remove the local folder
	err = os.RemoveAll(folderName)
	if err != nil {
		return err
	}
	return nil
}

// function to push branch to remote origin
func PushBranch(repo *git.Repository, branchName string) error {
	// push the branch to origin
	remote, err := repo.Remotes.Lookup("origin")

	cbs := &git.RemoteCallbacks{
		CredentialsCallback: credentialsCallback,
	}

	err = remote.Push([]string{"refs/heads/" + branchName}, &git.PushOptions{RemoteCallbacks: *cbs})
	if err != nil {
		log.Error("Error in Pushing the branch", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	return nil
}

// function to push branch to remote origin
func PushDeleteBranch(repo *git.Repository, branchName string) error {
	// push the branch to origin
	remote, err := repo.Remotes.Lookup("origin")

	cbs := &git.RemoteCallbacks{
		CredentialsCallback: credentialsCallback,
	}

	err = remote.Push([]string{":refs/heads/" + branchName}, &git.PushOptions{RemoteCallbacks: *cbs})
	if err != nil {
		log.Error("Error in Pushing the branch", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	return nil
}

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

	fmt.Println("Merge Message: ", mergeMessage)

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
	err := idx.RemoveByPath(fileName)
	if err != nil {
		return nil, err
	}

	err = deleteFile(path)
	if err != nil {
		return nil, err
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
func GitPull(folderName, branchName string) error {

	check, err := Exists(folderName)

	if !check {
		if err := os.Mkdir(folderName, os.ModePerm); err != nil {
			return err
		}
		// // clone the repo
		_, err := git.Clone("https://github.com/chitti-intel/test-flux-v3", folderName, &git.CloneOptions{CheckoutBranch: branchName, CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutSafe}})
		if err != nil {
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
		fmt.Println("Error in looking up Origin")
		return err
	}

	// Fetch changes from remote
	if err := remote.Fetch([]string{}, nil, ""); err != nil {
		fmt.Println("Error in Fetching")
		return err
	}

	// check if a local branch exitst if not do a checkout
	localBranch, err := repo.LookupBranch(branchName, git.BranchLocal)
	// No local branch, lets create one
	if localBranch == nil || err != nil {
		err = CheckoutBranch(folderName, branchName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + branchName)
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
		fmt.Println("Error in settting the head")
		return err
	}

	// Get remote master
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + branchName)
	if err != nil {
		return err
	}

	// testBranch := "testing-branch"
	// remoteBranchTest, err := repo.LookupBranch("origin/"+testBranch, git.BranchRemote)
	// if err != nil {
	// 	return err
	// }
	// log.Info("Test Branch Info ", log.Fields{"remoteBranchTest": remoteBranchTest, "commit": remoteBranchTest.Target()})
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

	fmt.Println("INFORMATION:")
	fmt.Println("remoteBranch")
	fmt.Println(remoteBranch)
	fmt.Println("commit")
	fmt.Println(remoteBranch.Target())
	fmt.Println("mergeRemoteHead")
	fmt.Println(mergeHeads)
	log.Info("Git Pull Information ", log.Fields{"remoteBranch": remoteBranch, "remoteBranchID": remoteBranchID, "annotatedCommit": annotatedCommit, "analysis": analysis, "mergeHeads": mergeHeads})

	// log.Info("MergeAnalysisFastForward", log.Fields{"analysis": analysis})
	// // Fast-forward changes
	// // Get remote tree
	// remoteTree, err := repo.LookupTree(remoteBranchID)
	// if err != nil {
	// 	log.Info("Error in looking up remote tree", log.Fields{"err": err})
	// 	return err
	// }

	// // Checkout
	// if err := repo.CheckoutTree(remoteTree, nil); err != nil {
	// 	log.Info("Error in checking out tree", log.Fields{"err": err})
	// 	return err
	// }

	// branchRef, err := repo.References.Lookup("refs/heads/" + branchName)
	// if err != nil {
	// 	log.Info("Error in looking up ref", log.Fields{"err": err})
	// 	return err
	// }

	// // Point branch to the object
	// branchRef.SetTarget(remoteBranchID, "")
	// if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
	// 	log.Info("Error in setting head traget", log.Fields{"err": err})
	// 	return err
	// }
	head, err := repo.Head()
	if err != nil {
		return err
	}

	if analysis&git.MergeAnalysisUpToDate != 0 {
		log.Info("MergeAnalysisUpToDate", log.Fields{"analysis": analysis})
		return nil
	} else if analysis&git.MergeAnalysisNormal != 0 {
		log.Info("MergeAnalysisNormal", log.Fields{"analysis": analysis})

		// set head to point to the created branch
		err = repo.SetHead("refs/heads/" + branchName)
		if err != nil {
			log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
			fmt.Println("Error in settting the head")
			return err
		}
		// Just merge changes
		if err := repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
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

		repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)

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

func PullBranch(repoPath string, remoteName string, branchName string, name string, email string) error {

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return err
	}
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + branchName)
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
		return err
	}

	remote, err := repo.Remotes.Lookup(remoteName)
	if err != nil {
		return err
	}

	// called := false
	cbs := &git.RemoteCallbacks{
		CredentialsCallback: credentialsCallback,
	}

	err = remote.Fetch([]string{}, &git.FetchOptions{RemoteCallbacks: *cbs}, "")

	if err != nil {
		log.Error("Error in Fetch", log.Fields{"err": err})
		return err
	}

	remoteBranch, err := repo.References.Lookup("refs/remotes/" + remoteName + "/" + branchName)
	if err != nil {
		log.Error("Error obtaining remote Branch ref", log.Fields{"err": err, "remoteBranch": remoteBranch})
		return err
	}

	mergeRemoteHead, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		log.Error("Error in obtaining annotated commit", log.Fields{"err": err})
		return err
	}

	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = mergeRemoteHead
	fmt.Println("INFORMATION:")
	fmt.Println("remoteBranch")
	fmt.Println(remoteBranch)
	fmt.Println("commit")
	fmt.Println(remoteBranch.Target())
	fmt.Println("mergeRemoteHead")
	fmt.Println(mergeRemoteHead)
	log.Info("Git Pull Information New ", log.Fields{"remoteBranch": remoteBranch, "mergeRemoteHead": mergeRemoteHead, "mergeHeads": mergeHeads})

	if err = repo.Merge(mergeHeads, nil, nil); err != nil {
		log.Error("Error in Mreging", log.Fields{"err": err, "mergeHeads": mergeHeads})
		return err
	}

	// Check if the index has conflicts after the merge
	idx, err := repo.Index()
	if err != nil {
		return err
	}

	currentBranch, err := repo.Head()
	if err != nil {
		return err
	}

	localCommit, err := repo.LookupCommit(currentBranch.Target())
	if err != nil {
		return err
	}

	// If index has conflicts, read old tree into index and
	// return an error.
	if idx.HasConflicts() {

		repo.ResetToCommit(localCommit, git.ResetHard, &git.CheckoutOpts{})

		repo.StateCleanup()

		return errors.New("conflict")
	}

	// If everything looks fine, create a commit with the two parents
	treeID, err := idx.WriteTree()
	if err != nil {
		return err
	}

	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return err
	}

	remoteCommit, err := repo.LookupCommit(remoteBranch.Target())
	if err != nil {
		return err
	}

	sig := &git.Signature{Name: name, Email: email, When: time.Now()}
	_, err = repo.CreateCommit("HEAD", sig, sig, "merged", tree, localCommit, remoteCommit)
	if err != nil {
		return err
	}

	repo.StateCleanup()

	return nil
}

func CheckoutBranch(folderName, branchName string) error {
	repo, err := git.OpenRepository(folderName)
	if err != nil {
		return err
	}
	checkoutOpts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutAllowConflicts | git.CheckoutUseTheirs,
	}
	//Getting the reference for the remote branch
	// remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + branchName)
	// remoteBranch, err := repo.LookupBranch("origin/"+branchName, git.BranchRemote)
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/" + branchName)
	if err != nil {
		// log1.Print("Failed to find remote branch: " + branchName)
		log.Error("Failed to find remote branch: ", log.Fields{"branchName": branchName})
		return err
	}
	defer remoteBranch.Free()

	// Lookup for commit from remote branch
	commit, err := repo.LookupCommit(remoteBranch.Target())
	if err != nil {
		// log1.Print("Failed to find remote branch commit: " + branchName)
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
			// log1.Print("Failed to create local branch: " + branchName)
			log.Error("Failed to create local branch: ", log.Fields{"branchName": branchName})
			return err
		}

		// Setting upstream to origin branch
		err = localBranch.SetUpstream("origin/" + branchName)
		if err != nil {
			// log1.Print("Failed to create upstream to origin/" + branchName)
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
		// log1.Print("Failed to lookup for commit in local branch " + branchName)
		log.Error("Failed to lookup for commit in local branch ", log.Fields{"branchName": branchName})
		return err
	}
	defer localCommit.Free()

	tree, err := repo.LookupTree(localCommit.TreeId())
	if err != nil {
		// log1.Print("Failed to lookup for tree " + branchName)
		log.Error("Failed to lookup for tree", log.Fields{"branchName": branchName})
		return err
	}
	defer tree.Free()

	// Checkout the tree
	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		// log1.Print("Failed to checkout tree " + branchName)
		log.Error("Failed to checkout tree ", log.Fields{"branchName": branchName})
		return err
	}
	// // Setting the Head to point to our branch
	// repo.SetHead("refs/heads/" + branchName)
	return nil
}

//Function to get files in a path
func GetFilesInPath(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return files, err
	}

	return files, nil
}

//Function to get File Contents
func GetFileContent(filePath string) (string, error) {

	file, err := os.Open(filePath)
	fileInfo, err := file.Stat()
	if err != nil {
		// error handling
		return "", err
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		return "", errors.New("path is a directory")
	}
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
func GitPullCMD(folderName, branchName string) error {
	mutex.Lock()
	defer mutex.Unlock()
	// // // open a repo
	check, err := Exists(folderName)

	if !check {
		if err := os.Mkdir(folderName, os.ModePerm); err != nil {
			return err
		}
		// // clone the repo
		_, err := git.Clone("https://github.com/chitti-intel/test-flux-v3", folderName, &git.CloneOptions{CheckoutBranch: branchName, CheckoutOptions: git.CheckoutOptions{Strategy: git.CheckoutSafe}})
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
		fmt.Println("Error in looking up Origin")
		return err
	}

	// Fetch changes from remote
	if err := remote.Fetch([]string{}, nil, ""); err != nil {
		fmt.Println("Error in Fetching")
		return err
	}
	// // check if a local branch exist if not do a checkout
	localBranch, err := repo.LookupBranch(branchName, git.BranchLocal)
	// // No local branch, lets create one
	if localBranch == nil || err != nil {
		fmt.Println("Checking Out the branch")
		err = CheckoutBranchCMD(folderName, branchName)
		if err != nil {
			fmt.Println(err)
		}
	}
	// set head to point to the created branch
	err = repo.SetHead("refs/heads/" + branchName)
	if err != nil {
		log.Error("Error in settting the head", log.Fields{"err": err, "branchName": branchName})
		fmt.Println("Error in settting the head")
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
	fmt.Println("Checking Out Branch")
	_, err := git.OpenRepository(folderName)
	if err != nil {
		return err
	}

	// using git command
	cmd := exec.Command("git", "checkout", "-b", branchName, "origin/"+branchName)
	cmd.Dir = folderName
	err = cmd.Run()
	if err != nil {
		fmt.Println("Git checkout returned error")
		fmt.Println(err)
		return err
	}
	return nil
}
