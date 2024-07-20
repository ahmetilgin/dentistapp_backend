package http

import (
	"backend/country"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc country.UseCase, authUC gin.HandlerFunc) {
	h := NewHandler(uc)

	regions := router.Group("/country")
	{
		regions.POST("", authUC, h.Create)
		regions.GET("", h.Get)
	}
}
