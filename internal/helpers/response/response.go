package response

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, status int, message string, err interface{}) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}
