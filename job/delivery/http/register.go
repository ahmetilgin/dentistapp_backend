package http

import (
	"backend/job"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc job.UseCase) {
	h := NewHandler(uc)

	jobs := router.Group("/jobs")
	{
		jobs.POST("", h.Create)
		jobs.GET("", h.Get)
		jobs.DELETE("", h.Delete)
	}
}
