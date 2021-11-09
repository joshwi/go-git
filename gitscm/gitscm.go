package gitscm

import (
	"fmt"
	"log"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func (p Project) Clone(directory string) (Project, error) {

	response := fmt.Sprintf(`[ Function: Clone ] [ Directory: %v ] [ Status: Success ]`, directory)

	repo := &git.Repository{}

	repo, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      p.Url,
		Progress: os.Stdout,
	})

	p.Repo = repo

	tree, err := p.Repo.Worktree()

	p.WorkTree = tree

	if err != nil {
		response = fmt.Sprintf(`[ Function: Clone ] [ Directory: %v ] [ Status: Failed ] [ Error: %v ]`, directory, err)
		log.Println(response)
		return p, err
	}

	log.Println(response)

	return p, nil

}

func (p Project) Open(directory string) (Project, error) {

	response := fmt.Sprintf(`[ Function: Open ] [ Directory: %v ] [ Status: Success ]`, directory)

	repo := &git.Repository{}

	repo, err := git.PlainOpen(directory)

	p.Repo = repo

	tree, err := p.Repo.Worktree()

	p.WorkTree = tree

	if err != nil {
		response = fmt.Sprintf(`[ Function: Open ] [ Directory: %v ] [ Status: Failed ] [ Error: %v ]`, directory, err)
		log.Fatalf(response)
		return p, err
	}

	log.Println(response)

	return p, nil

}

func (p Project) Add() error {

	response := fmt.Sprintf(`[ Function: Add ] [ Repo: %v ] [ Status: Success ]`, p.Name)

	_, err := p.WorkTree.Add(".")

	if err != nil {
		response = fmt.Sprintf(`[ Function: Add ] [ Repo: %v ] [ Status: Failed ] [ Error: %v ]`, p.Name, err)
		log.Fatalf(response)
		return err
	}

	log.Println(response)

	return nil

}

func (p Project) Commit(message string) error {

	response := fmt.Sprintf(`[ Function: Commit ] [ Message: %v ] [ Status: Success ]`, message)

	commit, err := p.WorkTree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  p.User,
			Email: p.Email,
			When:  time.Now(),
		},
	})

	if err != nil {
		response = fmt.Sprintf(`[ Function: Commit ] [ Message: %v ] [ Status: Failed ] [ Error: %v ]`, message, err)
		log.Fatalf(response)
		return err
	}

	_, err = p.Repo.CommitObject(commit)

	if err != nil {
		response = fmt.Sprintf(`[ Function: Commit ] [ Message: %v ] [ Status: Failed ] [ Error: %v ]`, message, err)
		log.Fatalf(response)
		return err
	}

	log.Println(response)

	return nil

}

func (p Project) Merge() error {

	// response := fmt.Sprintf(`[ Function: Merge ] [ Repo: %v ] [ Status: Success ]`, p.Name)

	return nil

}

func (p Project) Pull() error {

	// response := fmt.Sprintf(`[ Function: Pull ] [ Repo: %v ] [ Status: Success ]`, p.Name)

	// err := p.Repo.Push(&git.PushOptions{
	// 	Auth: &http.BasicAuth{
	// 		Username: p.User,
	// 		Password: p.Token,
	// 	},
	// })

	// if err != nil {
	// 	response = fmt.Sprintf(`[ Function: Pull ] [ Repo: %v ] [ Status: Failed ] [ Error: %v ]`, p.Name, err)
	// 	log.Fatalf(response)
	// 	return err
	// }

	// log.Println(response)

	return nil

}

func (p Project) Push() error {

	response := fmt.Sprintf(`[ Function: Push ] [ Repo: %v ] [ Status: Success ]`, p.Name)

	err := p.Repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: p.User,
			Password: p.Token,
		},
	})

	if err != nil {
		response = fmt.Sprintf(`[ Function: Push ] [ Repo: %v ] [ Status: Failed ] [ Error: %v ]`, p.Name, err)
		log.Fatalf(response)
		return err
	}

	log.Println(response)

	return nil

}

func (p Project) Branches(name string) ([]string, error) {

	output := []string{}

	branches, err := p.Repo.Branches()

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

func (p Project) Branch(name string, branch string) error {

	branch_name := fmt.Sprintf("refs/heads/%v", branch)

	headRef, err := p.Repo.Head()

	ref := plumbing.NewHashReference(plumbing.ReferenceName(branch_name), headRef.Hash())

	err = p.Repo.Storer.SetReference(ref)

	if err != nil {
		response := fmt.Sprintf(`[ Function: Branch ] [ Repo: %v ] [ Branch: %v ] [ Status: Failed ] [ Error: %v ]`, name, branch, err)
		log.Fatalf(response)
		return err
	}

	return nil

}
