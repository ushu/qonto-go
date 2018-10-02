package qonto

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
