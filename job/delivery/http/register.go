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
		jobs.GET("/search/:location/:keyword", h.Search)
		jobs.GET("/search_professions/:location/:profession", h.SearchProfession)
		jobs.GET("/get_populer_professions/:location", h.GetPopulerJobs)
		jobs.GET("", h.Get)
		jobs.DELETE("", authUC, h.Delete)
	}
}
