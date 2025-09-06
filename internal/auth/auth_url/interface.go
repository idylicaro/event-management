package auth_url

import (
	"github.com/gin-gonic/gin"
)

type GenerateAuthURLController interface {
	Handle(ctx *gin.Context)
}

type GenerateAuthURLService interface {
	Execute(providerName, userAgent, clientIP string) (string, error)
}
