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

func (s *TransactionService) CreateTransacton(ctx context.Context, request *models.Transaction) (models.CreateTransactionResponse, error) {
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
	err := s.TransactionRepository.CreateTransacton(ctx, request)
	if err != nil {
		return response, errors.Wrap(err, "failed to create transaction")
	}

	response.Reference = request.Reference
	response.TransactionStatus = request.TransactionStatus
	return response, nil

}

func (s *TransactionService) UpdateStatusTransacton(ctx context.Context, tokenData models.TokenData, request *models.UpdateTransactionStatus) error {

	//get transaction by reference
	dataTransaction, err := s.TransactionRepository.GetTransactonByReference(ctx, request.Reference, false)
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
	err = s.TransactionRepository.UpdateStatusTransacton(ctx, request.Reference, request.TransactionStatus, string(byteAdditionalInfo))
	if err != nil {
		return errors.Wrap(err, "failed to update transaction status")
	}
	return nil
}
