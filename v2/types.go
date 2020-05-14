package qonto

import (
	"time"
)

// Organization describes the holder of a Qonto account.
// Organisations can have several accounts attached.
type Organization struct {
	// Slug identifies the organization account
	Slug string `json:"slug"`
	// BankAccounts holds the list of all accounts for the Organization
	BankAccounts []*BankAccount `json:"bank_accounts"`
}

// BankAccount holds identification info about a single bank account.
// All fields are required and will be populated in the API response.
//
// To call the Transactions API, both "Slug" and "IBAN" values will be required.
type BankAccount struct {
	Slug                   string  `json:"slug"`                     // a unique identifier for the account
	IBAN                   string  `json:"iban"`                     // the IBAN (EU account unique ID)
	BIC                    string  `json:"BIC"`                      // the BIC (EU bank unique ID)
	Currency               string  `json:"currency"`                 // the currency for the account (usually "EUR")
	Balance                float64 `json:"balance"`                  // the account balance in EUR
	BalanceCents           int64   `json:"balance_cents"`            // the account balance in cents
	AuthorizedBalance      float64 `json:"authorized_balance"`       // the authorized balance in EUR
	AuthorizedBalanceCents int64   `json:"authorized_balance_cents"` // the authorized balance in cents
}

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
	ID                 string            `json:"transaction_id"`
	Amount             float64           `json:"amount"`
	AmountCents        int64             `json:"amount_cents"`
	LocalAmount        float64           `json:"local_amount"`
	LocalAmountCents   int64             `json:"local_amount_cents"`
	Side               TransactionSide   `json:"side"`
	OperationType      OperationType     `json:"operation_type"`
	Currency           string            `json:"currency"`
	LocalCurrency      string            `json:"local_currency"`
	SettledAt          *time.Time        `json:"settled_at,omitempty"`
	EmittedAt          time.Time         `json:"emitted_at"`
	Status             TransactionStatus `json:"status"`
	Note               *string           `json:"note,omitempty"`
	Label              *string           `json:"label,omitempty"`
	VATAmount          *float64          `json:"vat_amount,omitempty"`
	VATAmountCents     *int64            `json:"vat_amount_cents,omitempty"`
	VATRate            *float64          `json:"vat_rate,omitempty"`
	InitiatorID        *string           `json:"initiator_id,omitempty"`
	LabelIDs           []string          `json:"label_ids,omitempty"`
	AttachmentIDs      []string          `json:"attachment_ids,omitempty"`
	AttachmentLost     bool              `json:"attachment_lost,omitempty"`
	AttachmentRequired bool              `json:"attachment_required,omitempty"`
}

// Label as defined in the Qonto Dashboard
type Label struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

// Memberships as defined in the Qonto Dashboard
type Membership struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Attachment struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	FileName        string    `json:"file_name"`
	FileSize        int64     `json:"file_size,string,omitempty"`
	FileContentType string    `json:"file_content_type"`
	URL             string    `json:"url"`
}

// ResponseMeta holds the paging information when fetching a page of transactions.
type Meta struct {
	CurrentPage int  `json:"current_page"`
	NextPage    *int `json:"next_page"`
	PrevPage    *int `json:"prev_page"`
	TotalPages  int  `json:"total_pages"`
	TotalCount  int  `json:"total_count"`
	PerPage     int  `json:"per_page"`
}
