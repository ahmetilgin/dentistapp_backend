package http

import (
	"backend/job"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc job.UseCase, authUC gin.HandlerFunc ) {
	h := NewHandler(uc)
	jobs := router.Group("/jobs")
	{
		jobs.POST("", authUC, h.Create)
		jobs.GET("", h.Get)
		jobs.DELETE("",  authUC,h.Delete)
	}
}
