package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"orderly/internal/handler/http/middleware"
	"orderly/internal/pkg/token"
)

func initRoute(h handlerSet) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/signup", h.auth.Signup)
	router.POST("/login", h.auth.Login)

	protected := router.Group("", middleware.Auth(token.Parse))
	{
		protected.GET("/me", func(c *gin.Context) {
			userID, ok := middleware.UserIDFrom(c)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		})
	}

	return router
}
