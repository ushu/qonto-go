package qonto_test

import (
	"encoding/json"
	qonto "github.com/ushu/qonto-go"
	"testing"
)

var OrganizationFixture = `
{
  "slug": "test",
  "bank_accounts": [
    {
      "slug": "test-bank-account-1",
      "iban": "FR7600000000000000000000000",
      "bic": "TRZOFR21XXX",
      "currency": "EUR",
      "balance": 4000.36,
      "balance_cents": 400036,
      "authorized_balance": 4000.36,
      "authorized_balance_cents": 400036
    }
  ]
}
`

func TestOrganization_UnmarshalJSON(t *testing.T) {
	var o qonto.Organization
	err := json.Unmarshal([]byte(OrganizationFixture), &o)
	if err != nil {
		t.Fatalf("Could not parse JSON: %s", err.Error())
	}

	if o.Slug != "test" {
		t.Errorf("o.Slug == %q; want %q", o.Slug, "test")
	}
	if len(o.BankAccounts) != 1 {
		t.Errorf("len(BankAccounts) == %s; want %s", o.Slug, "test")
	}
}
