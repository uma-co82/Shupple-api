package controller

import (
	"github.com/gin-gonic/gin"
)

type HealthCheckController struct{}

func (healthCheckController HealthCheckController) HealthCheck(c *gin.Context) {
	c.JSON(200, "ok")
}
