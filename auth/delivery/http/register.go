package http

import (
	"backend/auth"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.Engine, uc auth.UseCase) {
	h := NewHandler(uc)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up-business-user", h.SignUpBusinessUser)
		authEndpoints.POST("/sign-up-normal-user", h.SignUpNormalUser)
		authEndpoints.POST("/sign-in-business-user", h.SignInBusinessUser)
		authEndpoints.POST("/sign-in-normal-user", h.SignInNormalUser)
	}
}
