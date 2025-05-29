package interfaces

import "github.com/gin-gonic/gin"

type IHealthCheckServices interface {
	HealthCheckServices() (string, error)
}

type IHealthCheckRepository interface {
}

type IHealthCheckApi interface {
	HealthChecHandlerHTTP(c *gin.Context)
}
