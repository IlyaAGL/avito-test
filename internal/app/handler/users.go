package handler

import (
	"context"

	"github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/users"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	SetUserActive(ctx context.Context, req users.SetActiveRequest) (*users.SetActiveResponse, error)
	GetUserReviewPRs(ctx context.Context, userID string) (*users.ReviewResponse, error)
	GetUser(ctx context.Context, userID string) (*users.UserResponse, error)
}

type userHandler struct {
	BaseHandler
	userService UserService
}

func NewUserHandler(userService UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) SetIsActive(c *gin.Context) {
	var req users.SetActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.BadRequest(c, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.UserID == "" {
		h.BadRequest(c, "INVALID_REQUEST", "UserID is required")
		return
	}

	response, err := h.userService.SetUserActive(c.Request.Context(), req)
	if err != nil {
		h.NotFound(c, "NOT_FOUND", "resource not found")
		return
	}

	h.Success(c, response)
}

func (h *userHandler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		h.BadRequest(c, "INVALID_REQUEST", "user_id is required")
		return
	}

	response, err := h.userService.GetUserReviewPRs(c.Request.Context(), userID)
	if err != nil {
		h.NotFound(c, "NOT_FOUND", "resource not found")
		return
	}

	h.Success(c, response)
}
