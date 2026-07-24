package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"orderly/internal/handler/http/response"
)

type Usecase interface {
	Signup(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type Handler struct {
	usecase Usecase
}

func New(usecase Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Signup(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "a valid email and a password are required")
		return
	}

	tok, err := h.usecase.Signup(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, tokenResponse{Token: tok})
}

func (h *Handler) Login(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "a valid email and a password are required")
		return
	}

	tok, err := h.usecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{Token: tok})
}
