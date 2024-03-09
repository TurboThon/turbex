package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/turbex-backend/src/handlers"
)

// @BasePath /api/v1
// @Summary API health check
// @Schemes
// @Description Performs a health check of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {string} API OK
// @Router /health [get]
func healthRoute(c *gin.Context) {
	handlers.DoStatus(c)
}
