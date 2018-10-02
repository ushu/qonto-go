package qonto

// Organization describes an Organization known to Qonto.
// It can have several accounts attached.
type Organization struct {
	Slug         string         `json:"slug"`          // a unique identifier of the organization (optional)
	BankAccounts []*BankAccount `json:"bank_accounts"` // all the accounts attached to this organization
}
