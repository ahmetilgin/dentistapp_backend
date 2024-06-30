package http

import (
	"backend/auth"
	"backend/job"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)


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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.CreateJob(c.Request.Context(), user, inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type getResponse struct {
	Jobs []*models.Job `json:"jobs"`
}

func (h *Handler) Get(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	bms, err := h.useCase.GetJobs(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getResponse{
		Jobs: toJobs(bms),
	})
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

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	if err := h.useCase.DeleteJob(c.Request.Context(), user, inp.ID); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func toJobs(bs []*models.Job) []*models.Job {
	out := make([]*models.Job, len(bs))

	for i, b := range bs {
		out[i] = toJob(b)
	}

	return out
}

func toJob(b *models.Job) *models.Job {
	return &models.Job{
		ID:     b.ID,
		UserID: b.UserID,
		JobTitle : b.JobTitle,
		Description : b.Description,
		Requirements : b.Requirements,
		Location : b.Location,
		SalaryRange : b.SalaryRange,
		EmploymentType: b.EmploymentType,  // full-time, part-time, contract, etc.
		DatePosted :  b.DatePosted,
		ApplicationDeadline :  b.ApplicationDeadline,
	}
}

