package auth_url

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type Controller struct {
	Service *Service
}

// @Summary Get the authentication URL for a provider
// @Description Get the authentication URL for a provider
// @Tags Auth
// @Accept json
// @Produce json
// @Param provider path string true "Provider name"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /auth/{provider} [get]
func (c *Controller) GetAuthURL(ctx *gin.Context) {
	provider := ctx.Param("provider") // Ex: "google" ou "github"
	url, err := c.Service.GenerateAuthURL(provider)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "auth.get_auth_url.failed", err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "auth.get_auth_url.success", url, nil)
}
