package http

import (
	"backend/region"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc region.UseCase,  authUC gin.HandlerFunc) {
	h := NewHandler(uc)

	regions := router.Group("/region")
	{
		regions.POST("", authUC,h.Create)
		regions.GET("", h.Get)
	}
}
