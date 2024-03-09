package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1
// @Summary API health check
// @Schemes
// @Description Performs a health check of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {string} API OK
// @Router /api/v1/status [get]
func DoStatus(c *gin.Context) {
  c.String(http.StatusOK, "API OK")
}

