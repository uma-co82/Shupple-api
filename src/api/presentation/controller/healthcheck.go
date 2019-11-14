package controller

import (
	"github.com/gin-gonic/gin"
)

type HealthCheckController struct{}

/**
 * ヘルスチェックに対して200を返す
 */
func (healthCheckController HealthCheckController) HealthCheck(c *gin.Context) {
	c.JSON(200, "ok")
}
