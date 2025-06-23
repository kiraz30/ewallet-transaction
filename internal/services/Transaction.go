package services

import (
	"context"
	"encoding/json"
	"ewallet-transaction/constants"
	"ewallet-transaction/external"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/models"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type TransactionService struct {
	TransactionRepository interfaces.ITransactionRepository
	External              interfaces.IExternal
}

func (s *TransactionService) CreateTransaction(ctx context.Context, request *models.Transaction) (models.CreateTransactionResponse, error) {
	var response models.CreateTransactionResponse

	request.TransactionStatus = constants.TransactionStatusPending
	request.Reference = helpers.GenerateReference()

	//validation json additional info
	jsonAdditionalInfo := map[string]interface{}{}
	if request.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(request.AdditionalInfo), &jsonAdditionalInfo)
		if err != nil {
			return response, errors.Wrap(err, "AdditionalInfo type is invalid")
		}
	}
	err := s.TransactionRepository.CreateTransaction(ctx, request)
	if err != nil {
		return response, errors.Wrap(err, "failed to create transaction")
	}

	response.Reference = request.Reference
	response.TransactionStatus = request.TransactionStatus
	return response, nil

}

func (s *TransactionService) UpdateStatusTransaction(ctx context.Context, tokenData models.TokenData, request *models.UpdateTransactionStatus) error {

	//get transaction by reference
	dataTransaction, err := s.TransactionRepository.GetTransactionByReference(ctx, request.Reference, false)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction by reference")
	}
	fmt.Println("additional info:", dataTransaction.AdditionalInfo)

	//validatio transaction status flow
	statusValid := false
	mapTransactionStatusFlow := constants.MapTransactionStatus[dataTransaction.TransactionStatus]
	for i := range mapTransactionStatusFlow {
		if mapTransactionStatusFlow[i] == request.TransactionStatus {
			statusValid = true
		}
	}

	if !statusValid {
		return fmt.Errorf("Transaction status flow invalid. request status = %s", request.TransactionStatus)
	}

	// request update balance to ewallet-wallet service
	requestUpdateBalance := external.UpdateBalance{
		Reference: request.Reference,
		Amount:    dataTransaction.Amount,
	}

	//for reversed
	if request.TransactionStatus == constants.TransactionStatusReversed {
		requestUpdateBalance.Reference = "REVERSED - " + request.Reference

		now := time.Now()
		expiredReversalTime := dataTransaction.CreatedAt.Add(constants.MaximumReversalDuration)
		if now.After(expiredReversalTime) {
			return errors.New("reversal time duration already expired")
		}
	}

	var (
		errUpdateBalance error
	)
	switch dataTransaction.TransactionType {
	case constants.TransactionTypeTopup:
		if request.TransactionStatus == constants.TransactionStatusSuccess {
			_, errUpdateBalance = s.External.CreditBalance(ctx, tokenData.Token, requestUpdateBalance)
		} else if request.TransactionStatus == constants.TransactionStatusReversed {
			_, errUpdateBalance = s.External.DebitBalance(ctx, tokenData.Token, requestUpdateBalance)
		}
	case constants.TransactionTypePurchase:
		if request.TransactionStatus == constants.TransactionStatusSuccess {
			_, errUpdateBalance = s.External.DebitBalance(ctx, tokenData.Token, requestUpdateBalance)
		} else if request.TransactionStatus == constants.TransactionStatusReversed {
			_, errUpdateBalance = s.External.CreditBalance(ctx, tokenData.Token, requestUpdateBalance)
		}
	}

	if errUpdateBalance != nil {
		errors.Wrap(errUpdateBalance, "failed to update balance")
	}

	//update additional info
	var (
		currentAdditionalInfo = map[string]interface{}{}
		newAdditionalInfo     = map[string]interface{}{}
	)

	if dataTransaction.AdditionalInfo != "" {
		err = json.Unmarshal([]byte(dataTransaction.AdditionalInfo), &currentAdditionalInfo)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal current additional info")
		}
	}
	fmt.Println("new additional info:", request.AdditionalInfo)

	if request.AdditionalInfo != "" {
		err = json.Unmarshal([]byte(request.AdditionalInfo), &newAdditionalInfo)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal new additional info")
		}
	}

	for key, val := range newAdditionalInfo {
		currentAdditionalInfo[key] = val
	}

	byteAdditionalInfo, err := json.Marshal(currentAdditionalInfo)
	if err != nil {
		return errors.Wrap(err, "failed to marshal merger additional info")
	}

	//update status in DB
	err = s.TransactionRepository.UpdateStatusTransaction(ctx, request.Reference, request.TransactionStatus, string(byteAdditionalInfo))
	if err != nil {
		return errors.Wrap(err, "failed to update transaction status")
	}
	dataTransaction.TransactionStatus = request.TransactionStatus
	s.sendNotification(ctx, tokenData, dataTransaction)
	return nil
}

func (s *TransactionService) sendNotification(ctx context.Context, tokenData models.TokenData, trx models.Transaction) {
	if trx.TransactionType == constants.TransactionTypePurchase && trx.TransactionStatus == constants.TransactionStatusSuccess {
		s.External.SendNotification(ctx, tokenData.Email, "purchase_success", map[string]string{
			"full_name":   tokenData.FullName,
			"description": trx.Description,
			"reference":   trx.Reference,
			"date":        trx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

}

func (s *TransactionService) GetTransaction(ctx context.Context, UserID string) ([]models.Transaction, error) {
	return s.TransactionRepository.GetTransaction(ctx, UserID)
}

func (s *TransactionService) GetTransactionDetail(ctx context.Context, reference string) (models.Transaction, error) {
	return s.TransactionRepository.GetTransactionByReference(ctx, reference, true)
}

func (s *TransactionService) RefundTransaction(ctx context.Context, tokenData models.TokenData, request *models.RefundTransaction) (models.CreateTransactionResponse, error) {
	var response models.CreateTransactionResponse

	dataTransaction, err := s.TransactionRepository.GetTransactionByReference(ctx, request.Reference, false)
	if err != nil {
		return response, errors.Wrap(err, "failed to get transaction by reference")
	}

	if dataTransaction.TransactionStatus != constants.TransactionStatusSuccess && dataTransaction.TransactionStatus != constants.TransactionTypePurchase {
		return response, errors.New("transaction status is not success or transaction type not purchase,  cannot refund")
	}
	referenceRefund := "REFUND - " + request.Reference
	requestUpdateBalance := external.UpdateBalance{
		Reference: referenceRefund,
		Amount:    dataTransaction.Amount,
	}

	_, err = s.External.CreditBalance(ctx, tokenData.Token, requestUpdateBalance)
	if err != nil {
		return response, errors.Wrap(err, "failed to credit balance")
	}
	transaction := models.Transaction{
		Amount:            dataTransaction.Amount,
		Reference:         referenceRefund,
		TransactionType:   constants.TransactionTypeRefund,
		TransactionStatus: constants.TransactionStatusSuccess,
		Description:       request.Description,
		AdditionalInfo:    request.AdditionalInfo,
	}

	err = s.TransactionRepository.CreateTransaction(ctx, &transaction)
	if err != nil {
		return response, errors.Wrap(err, "failed to insert new transaction transaction")
	}

	response.Reference = transaction.Reference
	response.TransactionStatus = transaction.TransactionStatus

	return response, nil
}
