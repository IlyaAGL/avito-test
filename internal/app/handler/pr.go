package handler

import (
	"context"

	"github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/common"
	pullrequests "github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/prs"
	"github.com/gin-gonic/gin"
)

type PullRequestService interface {
	CreatePR(ctx context.Context, req pullrequests.CreateRequest) (*pullrequests.CreateResponse, error)
	MergePR(ctx context.Context, req pullrequests.MergeRequest) (*pullrequests.MergeResponse, error)
	ReassignReviewer(ctx context.Context, req pullrequests.ReassignRequest) (*pullrequests.ReassignResponse, error)
	GetPR(ctx context.Context, prID string) (*pullrequests.PullRequestResponse, error)
	GetStats(ctx context.Context) (*common.StatsResponse, error)
}

type pullRequestHandler struct {
	BaseHandler
	prService PullRequestService
}

func NewpullRequestHandler(prService PullRequestService) *pullRequestHandler {
	return &pullRequestHandler{
		prService: prService,
	}
}

func (h *pullRequestHandler) CreatePR(c *gin.Context) {
	var req pullrequests.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.PullRequestID == "" {
		h.BadRequest(c, "INVALID_REQUEST", "PullRequestID is required")
		return
	}

	if req.PullRequestName == "" {
		h.BadRequest(c, "INVALID_REQUEST", "PullRequestName is required")
		return
	}

	if req.AuthorID == "" {
		h.BadRequest(c, "INVALID_REQUEST", "AuthorID is required")
		return
	}

	response, err := h.prService.CreatePR(c.Request.Context(), req)
	if err != nil {
		switch {
		case err.Error() == "author not found":
			h.NotFound(c, "NOT_FOUND", "Author not found")
		case err.Error() == "team not found":
			h.NotFound(c, "NOT_FOUND", "Team not found")
		default:
			h.Conflict(c, "PR_EXISTS", "PR already exists")
		}
		return
	}

	h.Created(c, response)
}

func (h *pullRequestHandler) MergePR(c *gin.Context) {
	var req pullrequests.MergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.PullRequestID == "" {
		h.BadRequest(c, "INVALID_REQUEST", "PullRequestID is required")
		return
	}

	response, err := h.prService.MergePR(c.Request.Context(), req)
	if err != nil {
		h.NotFound(c, "NOT_FOUND", "PR not found")
		return
	}

	h.Success(c, response)
}

func (h *pullRequestHandler) ReassignReviewer(c *gin.Context) {
	var req pullrequests.ReassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.PullRequestID == "" {
		h.BadRequest(c, "INVALID_REQUEST", "PullRequestID is required")
		return
	}

	if req.OldUserID == "" {
		h.BadRequest(c, "INVALID_REQUEST", "OldUserID is required")
		return
	}

	response, err := h.prService.ReassignReviewer(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case "cannot reassign on merged PR":
			h.Conflict(c, "PR_MERGED", "Cannot reassign on merged PR")
		case "reviewer is not assigned to this PR":
			h.Conflict(c, "NOT_ASSIGNED", "Reviewer is not assigned to this PR")
		case "no active replacement candidate in team":
			h.Conflict(c, "NO_CANDIDATE", "No active replacement candidate in team")
		default:
			h.NotFound(c, "NOT_FOUND", "PR or user not found")
		}
		return
	}

	h.Success(c, response)
}

func (h *pullRequestHandler) GetStats(c *gin.Context) {
	response, err := h.prService.GetStats(c.Request.Context())
	if err != nil {
		h.InternalError(c, "Failed to get statistics")
		return
	}

	h.Success(c, response)
}
