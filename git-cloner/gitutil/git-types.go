package gitutil

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitCloner interface {
	Clone(repoURL string, authType string) (bool, error)
}

type CIGitCloner struct {
	AuthType              string
	DefaultDirectory      string
	AllowBareRepositories bool
	WorkingBranch         string
	WorkTree              *git.Worktree
}

/*
1. Make random directory name


*/

func NewClient(workingBranch string) GitCloner {

	if workingBranch == "" {

		workingBranch = "main"
	}
	return &CIGitCloner{
		AuthType:              "",
		DefaultDirectory:      "/tmp/work",
		AllowBareRepositories: false,
		WorkingBranch:         workingBranch,
		WorkTree:              nil,
	}
}

func (cg *CIGitCloner) Clone(repoURL, authType string) (bool, error) {

	gitRepo, err := git.PlainClone(cg.DefaultDirectory, cg.AllowBareRepositories, &git.CloneOptions{
		URL:      "https://github.com/go-git/go-git",
		Progress: os.Stdout,
	})

	if err != nil {

		return false, err
	}

	wt, _ := gitRepo.Worktree()

	_ = gitRepo.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		fmt.Println(err)
	}
	wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", cg.WorkingBranch)),
		Force:  true,
	})

	ref, _ := gitRepo.Head()

	fInfo, err := wt.Filesystem.Stat("Dockerfile")

	if err != nil {

		return false, err
	}

	if fInfo.IsDir() {

		return false, fmt.Errorf("no dockerfile")

	}
	fmt.Println("Dockerfile exists. Size", fInfo.Size())

	fmt.Println("observed latest commit", ref.Hash())
	cg.WorkTree = wt
	return true, nil
}
