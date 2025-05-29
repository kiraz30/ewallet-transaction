package helpers

import "github.com/gin-gonic/gin"

type Response struct {
	Message string      `json:"message"`        // Message to be returned
	Data    interface{} `json:"data,omitempty"` // Data to be returned
}

func SendResponseHTTP(c *gin.Context, code int, message string, data interface{}) {
	resp := Response{
		Message: message,
		Data:    data,
	}
	c.JSON(code, resp)
}
