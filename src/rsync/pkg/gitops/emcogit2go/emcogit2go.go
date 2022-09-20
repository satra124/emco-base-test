package emcogit2go

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	log "gitlab.com/project-emco/core/emco-base/src/orchestrator/pkg/infra/logutils"

	git "github.com/libgit2/git2go/v33"
)

const (
	maxrand = 0x7fffffffffffffff
)

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

	// // // open a repo
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
	err = deleteBranch(repo, branchName)
	//push master to origin remote
	err = pushBranch(repo, "main")

	return nil
}

// function to push branch to remote origin
func pushBranch(repo *git.Repository, branchName string) error {
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
func deleteBranch(repo *git.Repository, branchName string) error {
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
