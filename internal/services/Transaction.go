package services

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/models"

	"github.com/pkg/errors"
)

type TransactionService struct {
	TransactionRepository interfaces.ITransactionRepository
}

func (s *TransactionService) CreateTransacton(ctx context.Context, request *models.Transaction) (models.CreateTransactionResponse, error) {
	var response models.CreateTransactionResponse

	request.TransactionStatus = constants.TransactionStatusPending
	request.Reference = helpers.GenerateReference()
	err := s.TransactionRepository.CreateTransacton(ctx, request)
	if err != nil {
		return response, errors.Wrap(err, "failed to create transaction")
	}

	response.Reference = request.Reference
	response.TransactionStatus = request.TransactionStatus
	return response, nil

}
