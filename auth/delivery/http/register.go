package http

import (
	"backend/auth"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc auth.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up-business", h.SignUpBusinessUser)
		authEndpoints.POST("/sign-up-normal", h.SignUpNormalUser)
		authEndpoints.POST("/sign-in", h.SignIn)
	}
}
