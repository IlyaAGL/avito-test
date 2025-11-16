package handler

import (
	"context"
	"strings"

	"github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/teams"
	"github.com/gin-gonic/gin"
)

type TeamService interface {
	CreateTeam(ctx context.Context, req teams.CreateRequest) (*teams.CreateResponse, error)
	GetTeam(ctx context.Context, teamName string) (*teams.TeamResponse, error)
	BulkDeactivateUsers(ctx context.Context, teamName string) error
}

type teamHandler struct {
	BaseHandler
	teamService TeamService
}

func NewTeamHandler(teamService TeamService) *teamHandler {
	return &teamHandler{
		teamService: teamService,
	}
}

func (h *teamHandler) AddTeam(c *gin.Context) {
	var req teams.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.TeamName == "" {
		h.BadRequest(c, "INVALID_REQUEST", "TeamName is required")
		return
	}

	if len(req.Members) == 0 {
		h.BadRequest(c, "INVALID_REQUEST", "Team must have at least one member")
		return
	}

	response, err := h.teamService.CreateTeam(c.Request.Context(), req)
	if err != nil {
		h.Conflict(c, "TEAM_EXISTS", "Team already exists")
		return
	}

	h.Created(c, response)
}

func (h *teamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		h.BadRequest(c, "INVALID_REQUEST", "team_name is required")
		return
	}

	response, err := h.teamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		h.NotFound(c, "NOT_FOUND", "resource not found")
		return
	}

	h.Success(c, response)
}

func (h *teamHandler) BulkDeactivateUsers(c *gin.Context) {
    var req teams.BulkDeactivateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.BadRequest(c, "INVALID_REQUEST", "Invalid request body")
        return
    }

    if req.TeamName == "" {
        h.BadRequest(c, "INVALID_REQUEST", "TeamName is required")
        return
    }

    err := h.teamService.BulkDeactivateUsers(c.Request.Context(), req.TeamName)
    if err != nil {
        if strings.Contains(err.Error(), "team not found") {
            h.NotFound(c, "NOT_FOUND", "Team not found")
        } else {
            h.InternalError(c, "Failed to deactivate users")
        }
        return
    }

    h.Success(c, gin.H{
        "message": "Deactivated",
        "team":    req.TeamName,
    })
}