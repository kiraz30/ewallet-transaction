package repository

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func (r *TransactionRepository) CreateTransacton(ctx context.Context, trx *models.Transaction) error {
	return r.DB.Create(&trx).Error

}

func (r *TransactionRepository) GetTransactonByReference(ctx context.Context, reference string, includeRefund bool) (models.Transaction, error) {
	var response models.Transaction
	sql := r.DB.Where("reference = ?", reference)
	//status refund cant update status
	if includeRefund {
		sql = sql.Where("transaction_status != ?", constants.TransactionTypeRefund)
	}

	err := sql.Last(&response).Error
	return response, err

}
func (r *TransactionRepository) UpdateStatusTransacton(ctx context.Context, reference, status, additional_info string) error {
	return r.DB.Exec("UPDATE transactions SET transaction_status = ? , additional_info = ? WHERE reference = ?", status, additional_info, reference).Error

}
