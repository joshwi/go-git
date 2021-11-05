package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/joshwi/go-utils/utils"
)

var (
	DIRECTORY = utils.Env("DIRECTORY")
	GIT_URL   = utils.Env("GIT_URL")
	GIT_USER  = utils.Env("GIT_USER")
	GIT_TOKEN = utils.Env("GIT_TOKEN")
)

func Clone(url string, directory string) (*git.Repository, error) {

	response := fmt.Sprintf(`[ Function: Clone ] [ Directory: %v ] [ Status: Success ]`, directory)

	repo := &git.Repository{}

	repo, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		response = fmt.Sprintf(`[ Function: Clone ] [ Directory: %v ] [ Status: Failed ] [ Error: %v ]`, directory, err)
		log.Println(response)
		return repo, err
	}

	log.Println(response)

	return repo, nil

}

func Open(directory string) (*git.Repository, error) {

	response := fmt.Sprintf(`[ Function: Open ] [ Directory: %v ] [ Status: Success ]`, directory)

	repo := &git.Repository{}

	repo, err := git.PlainOpen(directory)

	if err != nil {
		response = fmt.Sprintf(`[ Function: Open ] [ Directory: %v ] [ Status: Failed ] [ Error: %v ]`, directory, err)
		log.Println(response)
		return repo, err
	}

	log.Println(response)

	return repo, nil

}

func Add(tree *git.Worktree, directory string) (*git.Worktree, error) {

	filename := filepath.Join(directory, "md5.txt")
	response := fmt.Sprintf(`[ Function: Add ] [ Filename: %v ] [ Status: Success ]`, filename)
	err := ioutil.WriteFile(filename, []byte("1234567890"), 0644)

	if err != nil {
		response = fmt.Sprintf(`[ Function: Add ] [ Filename: %v ] [ Status: Failed ] [ Error: %v ]`, filename, err)
		log.Println(response)
		return nil, err
	}

	_, err = tree.Add("md5.txt")

	if err != nil {
		response = fmt.Sprintf(`[ Function: Add ] [ Filename: %v ] [ Status: Failed ] [ Error: %v ]`, filename, err)
		log.Println(response)
		return nil, err
	}

	log.Println(response)

	return tree, nil

}

func Commit(repo *git.Repository, tree *git.Worktree, message string) (*git.Worktree, error) {

	response := fmt.Sprintf(`[ Function: Commit ] [ Message: %v ] [ Status: Success ]`, message)

	commit, err := tree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "example",
			Email: "example@email.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		response = fmt.Sprintf(`[ Function: Commit ] [ Message: %v ] [ Status: Failed ] [ Error: %v ]`, message, err)
		log.Println(response)
		return nil, err
	}

	_, err = repo.CommitObject(commit)

	if err != nil {
		response = fmt.Sprintf(`[ Function: Commit ] [ Message: %v ] [ Status: Failed ] [ Error: %v ]`, message, err)
		log.Println(response)
		return nil, err
	}

	log.Println(response)

	return tree, nil

}

func Push(repo *git.Repository, name string) (*git.Repository, error) {

	response := fmt.Sprintf(`[ Function: Push ] [ Repo: %v ] [ Status: Success ]`, name)

	err := repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: GIT_USER,
			Password: GIT_TOKEN,
		},
	})

	if err != nil {
		response = fmt.Sprintf(`[ Function: Push ] [ Repo: %v ] [ Status: Failed ] [ Error: %v ]`, name, err)
		log.Println(response)
		return nil, err
	}

	log.Println(response)

	return repo, nil

}

func Branches(repo *git.Repository, name string) ([]string, error) {

	output := []string{}

	branches, err := repo.Branches()

	if err != nil {
		response := fmt.Sprintf(`[ Function: Branches ] [ Repo: %v ] [ Status: Failed ] [ Error: %v ]`, name, err)
		log.Println(response)
		return []string{}, nil
	}

	branches.ForEach(func(branch *plumbing.Reference) error {
		output = append(output, branch.Name().String())
		return nil
	})

	return output, nil

}

func Branch(repo *git.Repository, name string, branch string) (*git.Repository, error) {

	branch_name := fmt.Sprintf("refs/heads/%v", branch)

	headRef, err := repo.Head()

	ref := plumbing.NewHashReference(plumbing.ReferenceName(branch_name), headRef.Hash())

	// The created reference is saved in the storage.
	err = repo.Storer.SetReference(ref)

	if err != nil {
		response := fmt.Sprintf(`[ Function: Branch ] [ Repo: %v ] [ Branch: %v ] [ Status: Failed ] [ Error: %v ]`, name, branch, err)
		log.Println(response)
		return repo, err
	}

	return repo, nil

}

func main() {

	name := "nfldb-backup"

	dir := filepath.Join(DIRECTORY + "/" + name)

	repo, err := Clone(GIT_URL, dir)

	repo, err = Open(dir)

	if err != nil {
		log.Fatal(err)
	}

	branches, err := Branches(repo, name)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(branches)

	repo, err = Branch(repo, name, "v1")

	if err != nil {
		log.Fatal(err)
	}

	branches, err = Branches(repo, name)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(branches)

	// tree, err := repo.Worktree()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tree, err = Add(tree, dir)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tree, err = Commit(repo, tree, "test commit message")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// repo, err = Push(repo, name)

	// if err != nil {
	// 	log.Fatal(err)
	// }

}
