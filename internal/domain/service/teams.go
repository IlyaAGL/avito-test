package service

import (
	"context"
	"fmt"

	"github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/teams"
	"github.com/IlyaAGL/avito_autumn_2025/internal/models"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
	BulkDeactivateUsers(ctx context.Context, teamName string) error
}

type TeamService struct {
	teamRepo TeamRepository
}

func NewTeamService(teamRepo TeamRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, req teams.CreateRequest) (*teams.CreateResponse, error) {
	members := make([]models.Member, len(req.Members))
	for i, member := range req.Members {
		members[i] = models.Member{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	team := &models.Team{
		Name:    req.TeamName,
		Members: members,
	}

	err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return nil, err
	}

	memberResponses := make([]teams.MemberResponse, len(members))
	for i, member := range members {
		memberResponses[i] = teams.MemberResponse{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return &teams.CreateResponse{
		Team: teams.TeamResponse{
			TeamName: team.Name,
			Members:  memberResponses,
		},
	}, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*teams.TeamResponse, error) {
	team, err := s.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	memberResponses := make([]teams.MemberResponse, len(team.Members))
	for i, member := range team.Members {
		memberResponses[i] = teams.MemberResponse{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return &teams.TeamResponse{
		TeamName: team.Name,
		Members:  memberResponses,
	}, nil
}

func (s *TeamService) BulkDeactivateUsers(ctx context.Context, teamName string) error {
    _, err := s.teamRepo.GetTeam(ctx, teamName)
    if err != nil {
        return fmt.Errorf("team not found: %w", err)
    }

    err = s.teamRepo.BulkDeactivateUsers(ctx, teamName)
    if err != nil {
        return err
    }

    return nil
}