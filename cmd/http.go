package cmd

import (
	"ewallet-transaction/external"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/api"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	// d := dependencyInject()
	healthCheckSVC := &services.HealthCheck{}
	healtCheckAPI := &api.HealthCheck{
		HealthCheckServices: healthCheckSVC,
	}

	r := gin.Default()
	r.GET("/health", healtCheckAPI.HealthChecHandlerHTTP)
	err := r.Run(":" + helpers.GetEnv("PORT", "8081"))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	HealtyCheckApi interfaces.IHealthCheckApi
	External       interfaces.IExternal
}

func dependencyInject() Dependency {
	healtyCheckSVC := &services.HealthCheck{}
	healtyCheckAPI := &api.HealthCheck{
		HealthCheckServices: healtyCheckSVC,
	}
	external := &external.External{}

	return Dependency{
		HealtyCheckApi: healtyCheckAPI,
		External:       external,
	}

}
