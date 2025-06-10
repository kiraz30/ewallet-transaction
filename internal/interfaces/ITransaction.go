package interfaces

import (
	"context"
	"ewallet-transaction/internal/models"

	"github.com/gin-gonic/gin"
)

type ITransactionRepository interface {
	CreateTransacton(ctx context.Context, trx *models.Transaction) error
	GetTransactonByReference(ctx context.Context, reference string, includeRefund bool) (models.Transaction, error)
	UpdateStatusTransacton(ctx context.Context, reference, status, additional_info string) error
}

type ITransactionService interface {
	CreateTransacton(ctx context.Context, request *models.Transaction) (models.CreateTransactionResponse, error)
	UpdateStatusTransacton(ctx context.Context, tokenData models.TokenData, request *models.UpdateTransactionStatus) error
}

type ITransactionApi interface {
	CreateTransaction(c *gin.Context)
	UpdateStatusTransaction(c *gin.Context)
}
