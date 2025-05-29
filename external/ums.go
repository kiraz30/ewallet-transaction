package external

import (
	"context"
	"ewallet-transaction/constants"
	"ewallet-transaction/external/proto/tokenvalidation"
	"ewallet-transaction/internal/models"
	"fmt"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

type External struct {
}

func (e *External) ValidateToken(ctx context.Context, token string) (models.TokenData, error) {
	var (
		resp models.TokenData
	)

	conn, err := grpc.Dial("localhost:7000", grpc.WithInsecure())
	if err != nil {
		return resp, errors.Wrap(err, "faileded to connect to GRPC")
	}

	defer conn.Close()

	client := tokenvalidation.NewTokenValidationClient(conn)

	request := &tokenvalidation.TokenRequest{
		Token: token,
	}

	response, err := client.ValidateToken(ctx, request)
	if err != nil {
		return resp, errors.Wrap(err, "faileded to validate token")
	}

	if response.Message != constants.SuccessMessage {
		return resp, fmt.Errorf("got response error form UMS : %s", response.Message)

	}
	resp.UserID = response.Data.UserId
	resp.UserName = response.Data.Username
	resp.FullName = response.Data.FullName
	return resp, nil

}
