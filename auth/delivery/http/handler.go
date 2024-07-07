package http

import (
	"backend/auth"
	"backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) SignUpBusinessUser(c *gin.Context) {
	inp := new(models.BusinessUser)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.useCase.SignUpBusinessUser(c.Request.Context(), inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) SignUpNormalUser(c *gin.Context) {
	inp := new(models.NormalUser)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.useCase.SignUpNormalUser(c.Request.Context(), inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type signInResponse struct {
	Token string `json:"token"`
}

type signInResponseBusiness struct {
	Token string `json:"token"`
	BusinessName    string `json:"business_name"`
    BusinessAddress string `json:"business_address"`
}

type signInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignInNormalUser(c *gin.Context) {
	inp := new(signInput)
	fmt.Printf("Received JSON: %+v\n", c.Request.Body)

	if err := c.BindJSON(inp); err != nil {
		fmt.Printf("Error : %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, token, err := h.useCase.SignInNormalUser(c.Request.Context(), inp.Username, inp.Password)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, signInResponse{Token: token})
}

func (h *Handler) SignInBusinessUser(c *gin.Context) {
	inp := new(signInput)
	fmt.Printf("Received JSON: %+v\n", c.Request.Body)

	if err := c.BindJSON(inp); err != nil {
		fmt.Printf("Error : %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, token, err := h.useCase.SignInBusinessUser(c.Request.Context(), inp.Username, inp.Password)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, signInResponseBusiness{Token: token, BusinessName: user.BusinessName, BusinessAddress: user.BusinessAddress})
}
