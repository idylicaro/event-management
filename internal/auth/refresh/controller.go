package refresh

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/idylicaro/event-management/internal/dto/auth"
	"github.com/idylicaro/event-management/internal/helpers/response"
)

type refreshTokenController struct {
	Service RefreshTokenService
}

func NewRefreshTokenController(service RefreshTokenService) RefreshTokenController {
	return &refreshTokenController{Service: service}
}

func (c *refreshTokenController) Handle(context *gin.Context) {
	// Get refresh token from request body
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userAgent := context.GetHeader("User-Agent")
	clientIP := context.ClientIP()

	accessToken, refreshToken, err := c.Service.Execute(request.RefreshToken, userAgent, clientIP)
	if err != nil {
		response.Error(context, http.StatusUnauthorized, "Token refresh failed", err.Error())
		return
	}

	result := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response.Success(context, http.StatusOK, "Token refreshed successfully", result, nil)
}
