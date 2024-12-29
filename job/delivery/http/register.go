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
		jobs.POST("/search", h.Search)
		jobs.GET("/search_professions/:region/:profession", h.SearchProfession)
		jobs.GET("/get_populer_professions/:region", h.GetPopulerJobs)
		jobs.POST("/apply_job", authUC, h.ApplyJob)
		jobs.GET("/get_jobs", authUC, h.GetJobs)
		jobs.DELETE("/delete/:jobId", authUC, h.Delete)
		jobs.PUT("/update", authUC, h.Update)
		jobs.GET("/candidate/:candidateId", authUC, h.GetCandidateDetails)
	}
}
