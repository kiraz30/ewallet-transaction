package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Transaction struct {
	ID                int       `json:"-"`
	UserID            int       `json:"user_id"`
	Amount            float64   `json:"amount" gorm:"column:amount;type:decimal(15,2)" valid:"required"`
	TransactionType   string    `json:"transaction_type" gorm:"column:transaction_type;type:ENUM('TOPUP', 'PURCHASE', 'REFUND')"`
	TransactionStatus string    `json:"transaction_status" gorm:"column:transaction_status;type:ENUM('PENDING', 'SUCCESS', 'FAILED','REVERSED')" valid:"required"`
	Reference         string    `json:"reference" gorm:"column:reference;type:varchar(255)"`
	Description       string    `json:"description" gorm:"column:description;type:varchar(255)"`
	AdditionalInfo    string    `json:"additional_info" gorm:"column:additional_info";type:text`
	CreatedAt         time.Time `json:"-" `
	UpdatedAt         time.Time `json:"-" `
}

func (l *Transaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

func (*Transaction) TableName() string {
	return "transactions"
}

type CreateTransactionResponse struct {
	Reference         string `json:"reference"`
	TransactionStatus string `json:"transaction_status"`
}
