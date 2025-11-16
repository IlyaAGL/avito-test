package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"crypto/rand"
	"slices"
	"time"

	"github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/common"
	pullrequests "github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/prs"
	"github.com/IlyaAGL/avito_autumn_2025/internal/models"
)

type PullRequestRepository interface {
	CreatePR(ctx context.Context, pr *models.PullRequest) error
	GetPR(ctx context.Context, prID string) (*models.PullRequest, error)
	PRExists(ctx context.Context, prID string) (bool, error)
	MergePR(ctx context.Context, prID string) error
	UpdatePRReviewers(ctx context.Context, prID string, reviewerIDs []string) error
	GetReviewStats(ctx context.Context) ([]models.ReviewStats, error)
}

type pullRequestService struct {
	prRepo   PullRequestRepository
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewPullRequestService(prRepo PullRequestRepository, userRepo UserRepository, teamRepo TeamRepository) *pullRequestService {
	return &pullRequestService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *pullRequestService) CreatePR(ctx context.Context, req pullrequests.CreateRequest) (*pullrequests.CreateResponse, error) {
	author, err := s.userRepo.GetUser(ctx, req.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	teamMembers, err := s.userRepo.GetActiveTeamMembers(ctx, author.TeamName, []string{req.AuthorID})
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	reviewers := s.selectRandomReviewers(teamMembers, 2)

	reviewerIDs := make([]string, len(reviewers))
	for i, reviewer := range reviewers {
		reviewerIDs[i] = reviewer.UserID
	}

	pr := &models.PullRequest{
		ID:                req.PullRequestID,
		Name:              req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: reviewerIDs,
		CreatedAt:         time.Now(),
	}

	err = s.prRepo.CreatePR(ctx, pr)
	if err != nil {
		return nil, err
	}

	return &pullrequests.CreateResponse{
		PR: s.prToResponse(pr),
	}, nil
}

func (s *pullRequestService) MergePR(ctx context.Context, req pullrequests.MergeRequest) (*pullrequests.MergeResponse, error) {
	pr, err := s.prRepo.GetPR(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	if pr.Status == "MERGED" {
		return &pullrequests.MergeResponse{
			PR: s.prToResponse(pr),
		}, nil
	}

	err = s.prRepo.MergePR(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	mergedPR, err := s.prRepo.GetPR(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	return &pullrequests.MergeResponse{
		PR: s.prToResponse(mergedPR),
	}, nil
}

func (s *pullRequestService) ReassignReviewer(ctx context.Context, req pullrequests.ReassignRequest) (*pullrequests.ReassignResponse, error) {
	pr, err := s.prRepo.GetPR(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	if pr.Status == "MERGED" {
		return nil, errors.New("cannot reassign on merged PR")
	}

	if !slices.Contains(pr.AssignedReviewers, req.OldUserID) {
		return nil, errors.New("reviewer is not assigned to this PR")
	}

	oldReviewer, err := s.userRepo.GetUser(ctx, req.OldUserID)
	if err != nil {
		return nil, err
	}

	excludeIDs := pr.AssignedReviewers
	excludeIDs = append(excludeIDs, pr.AuthorID)

	teamMembers, err := s.userRepo.GetActiveTeamMembers(ctx, oldReviewer.TeamName, excludeIDs)
	if err != nil {
		return nil, err
	}

	if len(teamMembers) == 0 {
		return nil, errors.New("no active replacement candidate in team")
	}

	max := big.NewInt(int64(len(teamMembers)))

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}

	newReviewer := teamMembers[n.Int64()]

	newReviewers := make([]string, len(pr.AssignedReviewers))
	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == req.OldUserID {
			newReviewers[i] = newReviewer.UserID
		} else {
			newReviewers[i] = reviewer
		}
	}

	err = s.prRepo.UpdatePRReviewers(ctx, pr.ID, newReviewers)
	if err != nil {
		return nil, err
	}

	updatedPR, err := s.prRepo.GetPR(ctx, pr.ID)
	if err != nil {
		return nil, err
	}

	return &pullrequests.ReassignResponse{
		PR:         s.prToResponse(updatedPR),
		ReplacedBy: newReviewer.UserID,
	}, nil
}

func (s *pullRequestService) GetPR(ctx context.Context, prID string) (*pullrequests.PullRequestResponse, error) {
	pr, err := s.prRepo.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	response := s.prToResponse(pr)

	return &response, nil
}

func (s *pullRequestService) GetStats(ctx context.Context) (*common.StatsResponse, error) {
    reviewStats, err := s.prRepo.GetReviewStats(ctx)
	if err != nil {
		return nil, err
	}

	resultStats := make([]common.ReviewStats, 0, len(reviewStats))

	for _, stat := range reviewStats {
		resultStat := common.ReviewStats{
			UserID: stat.UserID,
			OpenPRs: stat.OpenPRs,
			TotalPRs: stat.TotalReviews,
		}

		resultStats = append(resultStats, resultStat)
	}

	response := &common.StatsResponse{
		UserStats: resultStats,
	}

	return response, nil
}

func (s *pullRequestService) selectRandomReviewers(users []models.User, max int) []models.User {
    if len(users) == 0 {
        return []models.User{}
    }

    shuffled := make([]models.User, len(users))
    copy(shuffled, users)

    for i := len(shuffled) - 1; i > 0; i-- {
        nBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
        if err != nil {
            return []models.User{}
        }

        j := int(nBig.Int64())

        shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
    }

    if len(shuffled) > max {
        return shuffled[:max]
    }

    return shuffled
}

func (s *pullRequestService) prToResponse(pr *models.PullRequest) pullrequests.PullRequestResponse {
	var mergedAtStr *string

	if pr.MergedAt != nil {
		mergedAt := pr.MergedAt.Format(time.RFC3339)
		mergedAtStr = &mergedAt
	}

	return pullrequests.PullRequestResponse{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
		MergedAt:          mergedAtStr,
	}
}
