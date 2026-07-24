package response

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"orderly/internal/entity"
)

type errorBody struct {
	Error string `json:"error"`
}

func Error(c *gin.Context, err error) {
	switch {
	case errors.Is(err, entity.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, errorBody{Error: err.Error()})
	case errors.Is(err, entity.ErrEmailTaken):
		c.JSON(http.StatusConflict, errorBody{Error: "email already registered"})
	case errors.Is(err, entity.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, errorBody{Error: "invalid credentials"})
	case errors.Is(err, entity.ErrUserNotFound):
		c.JSON(http.StatusNotFound, errorBody{Error: "not found"})
	default:
		log.Printf("unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, errorBody{Error: "internal server error"})
	}
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, errorBody{Error: message})
}
