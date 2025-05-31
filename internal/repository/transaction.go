package repository

import (
	"context"
	"ewallet-transaction/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func (r *TransactionRepository) CreateTransacton(ctx context.Context, trx *models.Transaction) error {
	return r.DB.Create(&trx).Error

}
