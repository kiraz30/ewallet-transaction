package constants

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
