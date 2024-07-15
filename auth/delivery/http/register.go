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
		authEndpoints.POST("/reset-password", h.ResetPassword)
		authEndpoints.POST("/send-email-normal-user", h.SendEmailNormalUser)
		authEndpoints.POST("/send-email-business-user", h.SendEmailBusinessUser)

	}
}
