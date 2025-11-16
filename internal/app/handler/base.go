package handler

import (
	"net/http"

	"github.com/IlyaAGL/avito_autumn_2025/internal/domain/dto/common"
	"github.com/gin-gonic/gin"
)

type BaseHandler struct{}

func (h *BaseHandler) Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func (h *BaseHandler) Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

func (h *BaseHandler) Error(c *gin.Context, code int, errorCode, message string) {
	c.JSON(code, common.ErrorResponse{
		Error: common.ErrorDetail{
			Code:    errorCode,
			Message: message,
		},
	})
}

func (h *BaseHandler) BadRequest(c *gin.Context, errorCode, message string) {
	h.Error(c, http.StatusBadRequest, errorCode, message)
}

func (h *BaseHandler) NotFound(c *gin.Context, errorCode, message string) {
	h.Error(c, http.StatusNotFound, errorCode, message)
}

func (h *BaseHandler) Conflict(c *gin.Context, errorCode, message string) {
	h.Error(c, http.StatusConflict, errorCode, message)
}

func (h *BaseHandler) InternalError(c *gin.Context, message string) {
	h.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
