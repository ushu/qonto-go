package qonto_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	qonto "github.com/ushu/qonto-go/v2"
)

func TestOrganization_FromJSON(t *testing.T) {
	j := loadFixture(t, "organization.json")

	var o qonto.Organization
	if err := json.Unmarshal([]byte(j), &o); err != nil {
		t.Fatalf("Could not parse JSON: %s", err.Error())
	}

	if o.Slug != "test-organization" {
		t.Errorf("o.Slug == %q; want %q", o.Slug, "test-organization")
	}
	if len(o.BankAccounts) != 1 {
		t.Errorf("len(BankAccounts) == %v; want %v", len(o.BankAccounts), 1)
	} else if o.BankAccounts[0].Slug != "test-bank-account" {
		t.Errorf("BankAccounts[0].Slug == %s; want %s", o.BankAccounts[0].Slug, "test-bank-account")
	}
}

func TestBankAccount_FromJSON(t *testing.T) {
	j := loadFixture(t, "bank_account.json")

	var b qonto.BankAccount
	if err := json.Unmarshal([]byte(j), &b); err != nil {
		t.Fatalf("Could not parse JSON: %s", err.Error())
	}

	if b.Slug != "test-bank-account" {
		t.Errorf("b.Slug == %q; want %q", b.Slug, "test-bank-account")
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

func TestTransaction_(t *testing.T) {
	j := loadFixture(t, "transaction.json")

	var tr qonto.Transaction
	if err := json.Unmarshal([]byte(j), &tr); err != nil {
		t.Fatalf("Could not parse JSON: %s", err.Error())
	}

	if tr.ID != "test-account-1-transaction-4" {
		t.Errorf("o.ID == %q; want %q", tr.ID, "test-account-1-transaction-4")
	}
	if tr.Amount != 120.42 {
		t.Errorf("o.Amount == %v; want %v", tr.Amount, 120.42)
	}
	if tr.AmountCents != 12042 {
		t.Errorf("o.AmountCents == %v; want %v", tr.AmountCents, 12042)
	}
	if tr.LocalAmount != 120.42 {
		t.Errorf("o.LocalAmount == %v; want %v", tr.LocalAmount, 120.42)
	}
	if tr.LocalAmountCents != 12042 {
		t.Errorf("o.LocalAmountCents == %v; want %v", tr.LocalAmountCents, 12042)
	}
	if tr.Side != qonto.TransactionSideDebit {
		t.Errorf("o.Side == %q; want %q", tr.Side, qonto.TransactionSideDebit)
	}
	if tr.OperationType != qonto.OperationTypeDirectDebit {
		t.Errorf("o.OperationType == %q; want %q", tr.OperationType, qonto.OperationTypeDirectDebit)
	}
	if tr.Currency != "EUR" {
		t.Errorf("o.Currency == %q; want %q", tr.Currency, "EUR")
	}
	if tr.LocalCurrency != "EUR" {
		t.Errorf("o.Currency == %q; want %q", tr.LocalCurrency, "EUR")
	}
	settledTime, _ := time.Parse(time.RFC3339, "2018-10-01T04:20:03.000Z")
	if *tr.SettledAt != settledTime {
		t.Errorf("o.SettledAt == %q; want %q", tr.SettledAt, settledTime)
	}
	emittedTime, _ := time.Parse(time.RFC3339, "2018-10-01T00:00:00.000Z")
	if tr.EmittedAt != emittedTime {
		t.Errorf("o.EmittedAt == %q; want %q", tr.EmittedAt, emittedTime)
	}
	if tr.Status != qonto.TransactionStatusCompleted {
		t.Errorf("o.Status == %q; want %q", tr.Status, qonto.TransactionStatusCompleted)
	}
	if tr.Note == nil {
		t.Errorf("o.Note == nil; want %q", "NOTE")
	} else if *tr.Note != "NOTE" {
		t.Errorf("o.Note == %q; want %q", *tr.Note, "NOTE")
	}
	if tr.Label == nil {
		t.Errorf("o.Label == nil; want %q", "SOME DEBIT")
	} else if *tr.Label != "SOME DEBIT" {
		t.Errorf("o.Label == %q; want %q", *tr.Label, "SOME DEBIT")
	}
}

func TestAttachment_FromJSON(t *testing.T) {
	j := loadFixture(t, "attachment.json")

	var a qonto.Attachment
	if err := json.Unmarshal([]byte(j), &a); err != nil {
		t.Fatalf("Could not parse JSON: %s", err.Error())
	}

	if a.ID != "1ec373a5-e30d-4a70-948d-c8d49e4a4d31" {
		t.Errorf("a.ID == %q; want %q", a.ID, "1ec373a5-e30d-4a70-948d-c8d49e4a4d31")
	}
	createdAt, _ := time.Parse(time.RFC3339, "2019-01-07T16:36:25.862Z")
	if a.CreatedAt != createdAt {
		t.Errorf("a.CreatedAt == %q; want %q", a.CreatedAt, createdAt)
	}
	if a.FileName != "doc.pdf" {
		t.Errorf("a.FileName == %q; want %q", a.FileName, "doc.pdf")
	}
	if a.FileSize != 49599 {
		t.Errorf("a.FileSize == %d; want %d", a.FileSize, 49599)
	}
	if a.FileContentType != "application/pdf" {
		t.Errorf("a.FileContentType == %q; want %q", a.FileContentType, "application/pdf")
	}
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()

	path := filepath.Join("testdata", name)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return f
}
