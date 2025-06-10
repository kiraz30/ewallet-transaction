package cmd

import (
	"ewallet-transaction/external"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/api"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/repository"
	"ewallet-transaction/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	d := dependencyInject()
	healthCheckSVC := &services.HealthCheck{}
	healtCheckAPI := &api.HealthCheck{
		HealthCheckServices: healthCheckSVC,
	}

	r := gin.Default()
	r.GET("/health", healtCheckAPI.HealthChecHandlerHTTP)

	transactionV1 := r.Group("/transaction/v1")
	transactionV1.POST("/create", d.MiddlewareValidateToken, d.TransactionApi.CreateTransaction)
	transactionV1.PUT("/update-status/:reference", d.MiddlewareValidateToken, d.TransactionApi.UpdateStatusTransaction)
	err := r.Run(":" + helpers.GetEnv("PORT", "8083"))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	HealtyCheckApi interfaces.IHealthCheckApi
	External       interfaces.IExternal
	TransactionApi interfaces.ITransactionApi
}

func dependencyInject() Dependency {
	healtyCheckSVC := &services.HealthCheck{}
	healtyCheckAPI := &api.HealthCheck{
		HealthCheckServices: healtyCheckSVC,
	}
	external := &external.External{}

	transactionRepository := &repository.TransactionRepository{
		DB: helpers.DB,
	}

	transactionService := &services.TransactionService{
		TransactionRepository: transactionRepository,
		External:              external,
	}

	transactionAPI := &api.TransactionAPI{
		TransactionService: transactionService,
	}

	return Dependency{
		HealtyCheckApi: healtyCheckAPI,
		External:       external,
		TransactionApi: transactionAPI,
	}

}
