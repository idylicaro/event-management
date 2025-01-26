package refresh

import "github.com/gin-gonic/gin"

type RefreshTokenController interface {
	Handle(ctx *gin.Context)
}

type RefreshTokenService interface {
	Execute(refreshToken string) (string, string, error)
}
