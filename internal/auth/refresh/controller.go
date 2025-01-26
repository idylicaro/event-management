package refresh

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type refreshTokenController struct {
	Service RefreshTokenService
}

func NewRefreshTokenController(service RefreshTokenService) RefreshTokenController {
	return &refreshTokenController{Service: service}
}

func (c *refreshTokenController) Handle(ctx *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, "auth.refreshToken.failed", err)
		return
	}

	// Refresh the token
	newAccessToken, newRefreshToken, err := c.Service.Execute(requestBody.RefreshToken)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "auth.refreshToken.failed", err)
		return
	}

	result := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	response.Success(ctx, http.StatusOK, "auth.refreshToken.success", result, nil)
}
