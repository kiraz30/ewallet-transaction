package api

import (
	"ewallet-transaction/constants"
	"ewallet-transaction/helpers"
	"ewallet-transaction/internal/interfaces"
	"ewallet-transaction/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionAPI struct {
	TransactionService interfaces.ITransactionService
}

func (api *TransactionAPI) CreateTransaction(c *gin.Context) {

	var (
		log     = helpers.Logger
		request models.Transaction
		// response models.CreateTransactionResponse
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("failed to get request :", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if err := request.Validate(); err != nil {
		log.Error("failed to validate request :", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		log.Error("failed to get token")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		log.Error("Failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	//valudation Transaction Type
	if !constants.MapTransactionType[request.TransactionType] {
		log.Error("invalid transaction type")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}
	request.UserID = int(tokenData.UserID)

	response, err := api.TransactionService.CreateTransacton(c.Request.Context(), &request)
	if err != nil {
		log.Error("failed to create transaction :", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, response)

}
