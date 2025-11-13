package errs

import (
	"errors"
	"fmt"
)

var (
	PullRequestMergedErr = errors.New("cannot reassign on merged PR")
	NotAssignedErr       = errors.New("reviewer is not assigned to this PR")
	NoCandidateErr       = errors.New("no active replacement candidate in team")
	NotFoundErr          = errors.New("resource not found")
)

type TeamExistsError struct {
	TeamName string
}

func (e TeamExistsError) Error() string {
	return fmt.Sprintf("%s already exists", e.TeamName)
}

type PullRequestExistsError struct {
	PullRequestID string
}

func (e PullRequestExistsError) Error() string {
	return fmt.Sprintf("%s already exists", e.PullRequestID)
}
