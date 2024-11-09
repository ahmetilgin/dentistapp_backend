package http

import (
	"backend/job"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc job.UseCase, authUC gin.HandlerFunc) {
	h := NewHandler(uc)
	jobs := router.Group("/jobs")
	{
		jobs.POST("", authUC, h.Create)
		jobs.GET("/search/:region/:location/:keyword", h.Search)
		jobs.GET("/search_professions/:region/:profession", h.SearchProfession)
		jobs.GET("/get_populer_professions/:region", h.GetPopulerJobs)
		jobs.DELETE("", authUC, h.Delete)
	}
}
