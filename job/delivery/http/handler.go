package http

import (
	"backend/auth"
	"backend/job"
	jobmongo "backend/job/repository/mongo"
	"backend/models"
	"net/http"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

type SearchOptions struct {
	Keyword  string `json:"keyword"`
	Location string `json:"location"`
	Region   string `json:"region"`
}

type SearchRequest struct {
	Position string `json:"position"`
	Region   string `json:"region"`
	Language string `json:"language"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet(auth.CtxUserKey)

	if err := h.useCase.CreateJob(c.Request.Context(), user.(*models.BusinessUser), inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Search(c *gin.Context) {
	inp := new(SearchRequest)
	if err := c.ShouldBindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if inp.Limit == 0 {
		inp.Limit = 10
	}
	if inp.Page == 0 {
		inp.Page = 1
	}

	result, err := h.useCase.Search(c.Request.Context(), inp.Position, inp.Region, inp.Page, inp.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	query := utils.NormalizeString(c.Param("profession"))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet(auth.CtxUserKey)

	if err := h.useCase.DeleteJob(c.Request.Context(), user.(*models.BusinessUser), inp.ID); err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet(auth.CtxUserKey)

	if err := h.useCase.ApplyJob(c.Request.Context(), user.(*models.NormalUser), inp.ID); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) GetJobs(c *gin.Context) {

	user := c.MustGet(auth.CtxUserKey)

	result, err := h.useCase.GetJobs(c.Request.Context(), user.(*models.BusinessUser))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &JobListResponse{
		Jobs: result,
	})
}

// Update handles the job update endpoint
func (h *Handler) Update(c *gin.Context) {
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet(auth.CtxUserKey)

	if err := h.useCase.Update(c.Request.Context(), user.(*models.BusinessUser), &job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "job updated successfully"})
}

func (h *Handler) GetCandidateDetails(c *gin.Context) {
	candidateID := c.Param("candidateId")
	if candidateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Candidate ID is required"})
		return
	}

	// user := c.MustGet(auth.CtxUserKey)
	// businessUser, ok := user.(*models.BusinessUser)
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }

	// // candidateDetails, err := h.useCase.GetCandidateDetails(c.Request.Context(), businessUser, candidateID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, candidateDetails)
}
