package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Erro interno no servidor",
				Error:   c.Errors.Last().Error(),
			})
		}
	}
}
