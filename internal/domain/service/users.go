package service

import (
	"context"

	"github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/users"
	"github.com/IlyaAGL/avito_autumn_2025/internal/models"
)

type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userID string) (*models.User, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserIDs []string) ([]models.User, error)
	GetUserReviewPRs(ctx context.Context, userID string) ([]models.PullRequestShort, error)
}

type userService struct {
	userRepo UserRepository
	prRepo   PullRequestRepository
}

func NewUserService(userRepo UserRepository, prRepo PullRequestRepository) *userService {
	return &userService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *userService) SetUserActive(ctx context.Context, req users.SetActiveRequest) (*users.SetActiveResponse, error) {
	user, err := s.userRepo.SetUserActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		return nil, err
	}

	return &users.SetActiveResponse{
		User: users.UserResponse{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}, nil
}

func (s *userService) GetUserReviewPRs(ctx context.Context, userID string) (*users.ReviewResponse, error) {
	_, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	prShorts, err := s.userRepo.GetUserReviewPRs(ctx, userID)
	if err != nil {
		return nil, err
	}

	prResponses := make([]users.PullRequestShort, len(prShorts))
	for i, pr := range prShorts {
		prResponses[i] = users.PullRequestShort{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		}
	}

	return &users.ReviewResponse{
		UserID:       userID,
		PullRequests: prResponses,
	}, nil
}

func (s *userService) GetUser(ctx context.Context, userID string) (*users.UserResponse, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &users.UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}, nil
}

