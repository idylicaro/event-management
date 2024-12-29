package callback

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type Controller struct {
	Service Service
}

func (c *Controller) HandleCallback(ctx *gin.Context) {
	provider := ctx.Param("provider")
	code := ctx.Query("code")

	user, err := c.Service.ProcessCallback(provider, code)
	if err != nil {
		response.Success(ctx, http.StatusBadRequest, "auth.callback.failed", err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "auth.callback.success", user, nil)
}
