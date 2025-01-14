package callback

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/idylicaro/event-management/internal/auth/providers"
	dto "github.com/idylicaro/event-management/internal/dto/auth"
	"github.com/idylicaro/event-management/internal/models"
)

type CallbackController interface {
	Handle(ctx *gin.Context)
}

type CallbackService interface {
	Execute(ctx context.Context, providerName, code string) (*dto.TokenResponse, error)
}

type CallbackRepository interface {
	FindOrCreateUser(data providers.UserInfo) (models.User, error)
}
