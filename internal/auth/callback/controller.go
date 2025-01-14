package callback

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type callbackController struct {
	Service CallbackService
}

func NewCallbackController(service CallbackService) CallbackController {
	return &callbackController{Service: service}
}

// @Summary Handle the callback from the authentication provider
// @Description Handle the callback from the authentication provider
// @Tags Auth
// @Accept json
// @Produce json
// @Param provider path string true "Provider name"
// @Param code query string true "Code from the provider"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /auth/{provider}/callback [get]
func (c *callbackController) Handle(ctx *gin.Context) {
	provider := ctx.Param("provider")
	code := ctx.Query("code")
	tokenResponse, err := c.Service.Execute(ctx, provider, code)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "auth.callback.failed", err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "auth.callback.success", tokenResponse, nil)
}
