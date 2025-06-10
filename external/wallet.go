package external

import (
	"bytes"
	"context"
	"encoding/json"
	"ewallet-transaction/helpers"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type UpdateBalance struct {
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}

type UpdateBalanceResponse struct {
	Messaage string `json:"message"`
	Data     struct {
		Balance float64 `json:"balance"`
	} `json:"data"`
}

func (e *External) CreditBalance(ctx context.Context, token string, req UpdateBalance) (*UpdateBalanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal json")
	}
	url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_CREDIT", "")
	fmt.Println("URL Wallet", url)
	httpRequest, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet http request")
	}

	httpRequest.Header.Set("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to wallet service")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from wallet service: %d", resp.StatusCode)
	}
	result := &UpdateBalanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	defer resp.Body.Close()
	return result, nil
}
func (e *External) DebitBalance(ctx context.Context, token string, req UpdateBalance) (*UpdateBalanceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal json")
	}
	url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_DEBIT", "")
	fmt.Println("URL Wallet", url)
	httpRequest, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet http reques")
	}

	httpRequest.Header.Set("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to wallet service")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from wallet service: %d", resp.StatusCode)
	}
	result := &UpdateBalanceResponse{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}
	defer resp.Body.Close()
	return result, nil
}
