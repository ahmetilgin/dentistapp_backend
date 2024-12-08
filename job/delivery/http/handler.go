package http

import (
	"backend/auth"
	"backend/job"
	jobmongo "backend/job/repository/mongo"
	"backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchOptions struct {
	Keyword  string `json:"keyword"`
	Location string `json:"location"`
	Region   string `json:"region"`
}
type Handler struct {
	useCase job.UseCase
}

func NewHandler(useCase job.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Create(c *gin.Context) {
	inp := new(models.Job)

	if err := c.BindJSON(inp); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.BusinessUser)

	if err := h.useCase.CreateJob(c.Request.Context(), user, inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Search(c *gin.Context) {
	inp := new(SearchOptions)
	location := c.Param("location")
	region := c.Param("region")
	keyword := c.Param("keyword")
	inp.Location = location
	inp.Keyword = keyword
	inp.Region = region
	result, err := h.useCase.Search(c.Request.Context(), inp.Location, inp.Keyword, inp.Region)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(result); i++ {
		result[i].JobDetail.Candidates = nil
	}

	c.JSON(http.StatusOK, &JobsResponse{
		Jobs: result,
	})
}

type queryResult struct {
	QueryResult []string `json:"query_result"`
}

func (h *Handler) SearchProfession(c *gin.Context) {

	query := c.Param("profession")
	location := c.Param("region")

	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "region parameter is required"})
		return
	}

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}
	result, err := h.useCase.SearchProfession(c.Request.Context(), query, location)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &queryResult{
		QueryResult: result,
	})
}

type JobsResponse struct {
	Jobs []*jobmongo.JobDetails `json:"jobs"`
}

type JobListResponse struct {
	Jobs []*models.Job `json:"jobs"`
}

type deleteInput struct {
	ID string `json:"id"`
}

func (h *Handler) Delete(c *gin.Context) {
	inp := new(deleteInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.BusinessUser)

	if err := h.useCase.DeleteJob(c.Request.Context(), user, inp.ID); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetPopulerJobs(c *gin.Context) {

	code := c.Param("region")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "region parameter is required"})
		return
	}

	result, err := h.useCase.GetPopulerJobs(c.Request.Context(), code)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &queryResult{
		QueryResult: result,
	})
}

type applyJobInput struct {
	ID string `json:"job_id"`
}

func (h *Handler) ApplyJob(c *gin.Context) {
	inp := new(applyJobInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	jobId := inp.ID

	if jobId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := c.MustGet(auth.CtxUserKey).(*models.NormalUser)
	if err := h.useCase.ApplyJob(c.Request.Context(), user, jobId); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) GetJobs(c *gin.Context) {

	user := c.MustGet(auth.CtxUserKey).(*models.BusinessUser)

	result, err := h.useCase.GetJobs(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &JobListResponse{
		Jobs: result,
	})
}
