package http

import (
	"backend/auth"
	"backend/models"
	"backend/region"
	"net/http"

	"github.com/gin-gonic/gin"
)


type Handler struct {
	useCase region.UseCase
}

func NewHandler(useCase region.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}


func (h *Handler) Create(c *gin.Context) {
	inp := new(models.Region)
 	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.MustGet(auth.CtxUserKey)
	
	if err := h.useCase.CreateRegion(c.Request.Context() ,inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type getResponse struct {
	QueryResult []string `json:"query_result"`
}

func (h *Handler) Get(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	bms, err := h.useCase.Search(c.Request.Context(), query)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &getResponse{
		QueryResult: bms,
	})
}



