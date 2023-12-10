package logger

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendErrorsToClient(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		c.JSON(http.StatusBadRequest, c.Errors)
	}
}
