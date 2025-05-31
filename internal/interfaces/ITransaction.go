package interfaces

import (
	"context"
	"ewallet-transaction/internal/models"

	"github.com/gin-gonic/gin"
)

type ITransactionRepository interface {
	CreateTransacton(ctx context.Context, trx *models.Transaction) error
}

type ITransactionService interface {
	CreateTransacton(ctx context.Context, request *models.Transaction) (models.CreateTransactionResponse, error)
}

type ITransactionApi interface {
	CreateTransaction(c *gin.Context)
}
