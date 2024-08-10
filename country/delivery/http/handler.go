package http

import (
	"backend/auth"
	"backend/country"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase country.UseCase
}

func NewHandler(useCase country.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Create(c *gin.Context) {
	inp := new(models.Country)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.MustGet(auth.CtxUserKey)

	if err := h.useCase.CreateRegion(c.Request.Context(), inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type getResponse struct {
	QueryResult []string `json:"query_result"`
}

func (h *Handler) Get(c *gin.Context) {
	query := c.Param("keyword")
	code := c.Param("region")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code parameter is required"})
		return
	}

	bms, err := h.useCase.Search(c.Request.Context(), query, code)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getResponse{
		QueryResult: bms,
	})
}
