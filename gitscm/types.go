package gitscm

import git "github.com/go-git/go-git/v5"

type Project struct {
	Name      string
	Directory string
	Url       string
	Token     string
	User      string
	Email     string
	Repo      *git.Repository
	WorkTree  *git.Worktree
}
