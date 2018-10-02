package qonto_test

import (
	"encoding/json"
	qonto "github.com/ushu/qonto-go"
	"testing"
)

var BankAccountFixture = `
{
  "slug": "test-bank-account-1",
  "iban": "FR7600000000000000000000000",
  "bic": "TRZOFR21XXX",
  "currency": "EUR",
  "balance": 4000.36,
  "balance_cents": 400036,
  "authorized_balance": 3000.36,
  "authorized_balance_cents": 300036
}
`

func TestBankAccount(t *testing.T) {
	var b qonto.BankAccount
	err := json.Unmarshal([]byte(BankAccountFixture), &b)
	if err != nil {
		t.Fatalf("Could not parse JSON: %s", err.Error())
	}

	if b.Slug != "test-bank-account-1" {
		t.Errorf("b.Slug == %q; want %q", b.Slug, "test-bank-account-1")
	}
	if b.IBAN != "FR7600000000000000000000000" {
		t.Errorf("b.IBAN == %q; want %q", b.IBAN, "FR7600000000000000000000000")
	}
	if b.BIC != "TRZOFR21XXX" {
		t.Errorf("b.BIC == %q; want %q", b.BIC, "TRZOFR21XXX")
	}
	if b.Currency != "EUR" {
		t.Errorf("b.Currency == %q; want %q", b.Currency, "EUR")
	}
	if b.Balance != 4000.36 {
		t.Errorf("b.Balance == %v; want %v", b.Balance, 4000.36)
	}
	if b.BalanceCents != 400036 {
		t.Errorf("b.BalanceCents == %v; want %v", b.BalanceCents, 400036)
	}
	if b.AuthorizedBalance != 3000.36 {
		t.Errorf("b.AuthorizedBalance == %v; want %v", b.AuthorizedBalance, 3000.36)
	}
	if b.AuthorizedBalanceCents != 300036 {
		t.Errorf("b.AuthorizedBalanceCents == %v; want %v", b.AuthorizedBalanceCents, 300036)
	}
}
