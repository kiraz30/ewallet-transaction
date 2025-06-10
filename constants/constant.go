package constants

import "time"

const (
	ErrFailedBadRequest = "data tidak sesuai"
	ErrServerError      = "terjadi kesalahan pada server"
	SuccessMessage      = "success"
)

const (
	TransactionStatusPending  = "PENDING"
	TransactionStatusSuccess  = "SUCCESS"
	TransactionStatusFailed   = "FAILED"
	TransactionStatusReversed = "REVERSED"
)

const (
	TransactionTypePurchase = "PURCHASE"
	TransactionTypeTopup    = "TOPUP"
	TransactionTypeRefund   = "REFUND"
)

var MapTransactionType = map[string]bool{
	TransactionTypePurchase: true,
	TransactionTypeTopup:    true,
	TransactionTypeRefund:   true,
}

// flow perubahan status transaksi
var MapTransactionStatus = map[string][]string{
	TransactionStatusPending: {TransactionStatusSuccess, TransactionStatusFailed},
	TransactionStatusSuccess: {TransactionStatusReversed},
	TransactionStatusFailed:  {TransactionStatusSuccess},
}

const (
	MaximumReversalDuration = time.Hour * 24 // 24 hours
)
