package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DoStatus(c *gin.Context) {
	c.String(http.StatusOK, "API OK")
}
