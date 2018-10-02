package qonto

import (
	"time"
)

// TransactionSide holds The side of a transaction (debit/credit)
type TransactionSide string

const (
	// TransactionSideDebit marks an outgoing transaction
	TransactionSideDebit TransactionSide = "debit"
	// TransactionSideCredit marks an incoming transaction
	TransactionSideCredit = "credit"
)

// OperationType holds the type of transaction, as stored by Qonto.
type OperationType string

const (
	// OperationTypeTransfer marks a transfer.
	OperationTypeTransfer OperationType = "transfer"
	// OperationTypeCard marks a credit card operation.
	OperationTypeCard = "card"
	// OperationTypeDirectDebit marks a debit.
	OperationTypeDirectDebit = "direct_debit"
	// OperationTypeDirectIncome marks a debit.
	OperationTypeDirectIncome = "income"
	// OperationTypeDirectQontoFee marks a bank fee.
	OperationTypeDirectQontoFee = "qonto_fee"
)

// TransactionStatus indicates the status (processing etc.) of a transaction.
type TransactionStatus string

const (
	// TransactionStatusPending indicates the transaction is currently being processed.
	TransactionStatusPending TransactionStatus = "pending"
	// TransactionStatusReversed indicates the transaction has been reversed.
	TransactionStatusReversed = "reversed"
	// TransactionStatusDeclined indicates the transaction payment has been declined.
	TransactionStatusDeclined = "declined"
	// TransactionStatusCompleted indicates the transaction is finished and reflected on the amount.
	TransactionStatusCompleted = "completed"
)

// Transaction holds the detail of a transaction in an Account.
type Transaction struct {
	ID               string            `json:"transaction_id"`     // unique ID for the transaction
	Amount           float64           `json:"amount"`             // amount in EUR
	AmountCents      int64             `json:"amount_cents"`       // amount in cents
	LocalAmount      float64           `json:"local_amount"`       // amount in the original currency
	LocalAmountCents int64             `json:"local_amount_cents"` // amount in cents in the original currency
	Side             TransactionSide   `json:"side"`               // "debit" or "credit"
	OperationType    OperationType     `json:"operation_type"`     // type of transaction
	Currency         string            `json:"currency"`           // account currency, usually EUR
	LocalCurrency    string            `json:"local_currency"`     // original transaction currenyc (can be ≠ EUR)
	SettledAt        time.Time         `json:"settled_at"`         // date the debit/credit impacted the amount
	EmittedAt        time.Time         `json:"emitted_at"`         // date the transaction actually happened
	Status           TransactionStatus `json:"status"`             // current status of the transaction
	Note             string            `json:"note"`               // (optional) note writter in the Qonto UI
	Label            string            `json:"label"`              // Original label from the payer
}
