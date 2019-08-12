package history

import (
	"errors"
	"log"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// ErrCommonCommitFound is used for identifying when the iterator has reached the common commit
var ErrCommonCommitFound = errors.New("common commit found")

// CommitsOnBranch iterates through all references and returns commit hashes on given branch
func CommitsOnBranch(
	repo *git.Repository,
	branchRef plumbing.ReferenceName,
	compareRef plumbing.ReferenceName,
) ([]plumbing.Hash, error) {

	var commits []plumbing.Hash

	// common ancestor between two branches based on MergeBase
	var commonHash plumbing.Hash

	branchCommit, err := commitFromRepo(repo, branchRef.String())

	if err != nil {
		return nil, err
	}

	compareCommit, err := commitFromRepo(repo, compareRef.String())

	if err != nil {
		return nil, err
	}

	diffCommits, mergeBaseErr := branchCommit.MergeBase(compareCommit)

	if mergeBaseErr != nil {
		log.Print(mergeBaseErr)
		return nil, mergeBaseErr
	}

	commonHash = diffCommits[len(diffCommits)-1].Hash

	commitIter, logErr := repo.Log(&git.LogOptions{
		From: branchCommit.Hash,
	})

	if logErr != nil {
		return nil, logErr
	}

	iterErr := commitIter.ForEach(func(commit *object.Commit) error {
		if commit.Hash == commonHash {
			return ErrCommonCommitFound
		}
		commits = append(commits, commit.Hash)
		return nil
	})

	if iterErr != nil {
		if iterErr == ErrCommonCommitFound {
			return commits, nil
		}
		return nil, iterErr
	}

	return commits, nil
}

func commitFromRepo(repo *git.Repository, desiredBranch string) (*object.Commit, error) {
	desiredHash, err := repo.ResolveRevision(plumbing.Revision(desiredBranch))

	if err != nil {
		return nil, err
	}

	desiredCommit, err := repo.CommitObject(*desiredHash)

	if err != nil {
		return nil, err
	}

	return desiredCommit, nil
}