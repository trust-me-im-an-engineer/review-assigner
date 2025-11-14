package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"review-assigner/internal/errs"
	"review-assigner/internal/model"
	"review-assigner/internal/storage"
	"review-assigner/internal/storage/transaction"
)

// Service is considered to be a core layer of witch only one could exist,
// therefore it doesn't use interface
type Service struct {
	storage   storage.Storage
	txManager transaction.Manager
}

func NewService(storage storage.Storage, txManager transaction.Manager) *Service {
	return &Service{
		storage:   storage,
		txManager: txManager,
	}
}

func (s *Service) AddTeamAddUpdateUsers(ctx context.Context, team *model.Team) (*model.Team, error) {
	var result *model.Team

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		teamName, err := s.storage.AddTeam(ctx, team.TeamName)
		if err != nil {
			return fmt.Errorf("storage failed to add team: %w", err)
		}

		inputUsers := make([]model.User, len(team.Members))
		for i, member := range team.Members {
			inputUsers[i] = model.User{
				UserID:   member.UserID,
				Username: member.Username,
				TeamName: team.TeamName,
				IsActive: member.IsActive,
			}
		}

		users, err := s.storage.AddUpdateUsers(ctx, inputUsers)
		if err != nil {
			return fmt.Errorf("storage failed to add/update users: %w", err)
		}

		// Strictly speaking storage doesn't add anything new to users,
		// so we could just return team parameter,
		// but for future-proofing and consistency sake we parse team back from storage

		members := make([]model.TeamMember, len(users))
		for i, user := range users {
			members[i] = model.TeamMember{
				UserID:   user.UserID,
				Username: user.Username,
				IsActive: user.IsActive,
			}
		}

		result = &model.Team{
			TeamName: teamName,
			Members:  members,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GetTeam(ctx context.Context, name string) (*model.Team, error) {
	team, err := s.storage.GetTeam(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("storage failed to get team: %w", err)
	}
	return team, nil
}

func (s *Service) SetUserActivity(ctx context.Context, id string, active bool) (*model.User, error) {
	user, err := s.storage.SetUserActivity(ctx, id, active)
	if err != nil {
		return nil, fmt.Errorf("storage failed to set user activity: %w", err)
	}
	return user, nil
}

// CreatePullRequest ignores status field in model.PullRequestShort
func (s *Service) CreatePullRequest(ctx context.Context, pr *model.PullRequestShort) (*model.PullRequest, error) {
	var result *model.PullRequest

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		activeColleges, err := s.storage.GetActiveColleges(ctx, pr.AuthorID)
		if err != nil {
			return fmt.Errorf("storage failed to get active collegs: %w", err)
		}

		// Pick two distinct reviewers at random
		reviewers := make([]string, 0, 2)
		if len(activeColleges) <= 2 {
			reviewers = activeColleges
		} else {
			i1 := rand.Int() % len(activeColleges)
			reviewers = append(reviewers, activeColleges[i1])
			i2 := rand.Int() % (len(activeColleges) - 1)
			// increment to avoid picking same college twice
			if i2 >= i1 {
				i2++
			}
			reviewers = append(reviewers, activeColleges[i2])
		}

		createdAt := time.Now()
		inputPR := &model.PullRequest{
			PullRequestID:     pr.PullRequestID,
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID,
			Status:            model.PullRequestStatusOPEN,
			AssignedReviewers: reviewers,
			CreatedAt:         &createdAt,
			MergedAt:          nil,
		}

		result, err = s.storage.CreatePullRequest(ctx, inputPR)
		if err != nil {
			return fmt.Errorf("storage failed to create pull request: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) MergePullRequest(ctx context.Context, id string) (*model.PullRequest, error) {
	var result *model.PullRequest

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		pr, err := s.storage.GetPullRequest(ctx, id)
		if err != nil {
			return fmt.Errorf("storage failed to get pull request: %w", err)
		}

		pr.Status = model.PullRequestStatusMERGED
		mergedAt := time.Now()
		pr.MergedAt = &mergedAt

		result, err = s.storage.UpdatePullRequest(ctx, pr)
		if err != nil {
			return fmt.Errorf("storage failed to update pull request: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) ReassignPullRequest(ctx context.Context, pullRequestID, oldReviewerID string) (pr *model.PullRequest, newReviewerID string, err error) {
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		pr, err := s.storage.GetPullRequest(ctx, pullRequestID)
		if err != nil {
			return fmt.Errorf("storage failed to get pull request: %w", err)
		}

		if pr.Status == model.PullRequestStatusMERGED {
			return errs.PullRequestMergedErr
		}

		if !slices.Contains(pr.AssignedReviewers, oldReviewerID) {
			return errs.NotAssignedErr
		}

		// Search by author id and not oldReviewerID because reviewer could've changed team,
		// and we need original team to review pr.
		activeColleges, err := s.storage.GetActiveColleges(ctx, pr.AuthorID)
		if err != nil {
			return fmt.Errorf("storage failed to get active colleges: %w", err)
		}

		// remove oldReviewer to avoid setting same reviewer
		if i := slices.Index(activeColleges, oldReviewerID); i != -1 {
			activeColleges = slices.Delete(activeColleges, i, i+1)
		}

		if len(activeColleges) < 1 {
			return errs.NoCandidateErr
		}

		newReviewerID = activeColleges[rand.Int()%len(activeColleges)]

		i := slices.Index(pr.AssignedReviewers, oldReviewerID)
		pr.AssignedReviewers[i] = newReviewerID

		pr, err = s.storage.UpdatePullRequest(ctx, pr)
		if err != nil {
			return fmt.Errorf("stoage failed to update pull request: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, "", err
	}
	return pr, newReviewerID, nil
}

func (s *Service) GetUserAssignments(ctx context.Context, id string) ([]model.PullRequestShort, error) {
	prs, err := s.storage.GetUserAssignments(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("storage failed to get user assignments: %w", err)
	}

	shortPRs := make([]model.PullRequestShort, len(prs))
	for i, pr := range prs {
		shortPRs[i] = model.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		}
	}

	return shortPRs, nil
}
