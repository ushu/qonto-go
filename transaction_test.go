package qonto_test

import (
	"encoding/json"
	"testing"
	"time"

	qonto "github.com/ushu/qonto-go"
)

const TransactionFixture = `
{
  "transaction_id": "test-account-1-transaction-4",
  "amount": 120.42,
  "amount_cents": 12042,
  "local_amount": 120.42,
  "local_amount_cents": 12042,
  "side": "debit",
  "operation_type": "direct_debit",
  "currency": "EUR",
  "local_currency": "EUR",
  "label": "SOME DEBIT",
  "settled_at": "2018-10-01T04:20:03.000Z",
  "emitted_at": "2018-10-01T00:00:00.000Z",
  "status": "completed",
  "note": "NOTE"
}
`

func TestTransaction_UnmarshalJSON(t *testing.T) {
	var tr qonto.Transaction
	err := json.Unmarshal([]byte(TransactionFixture), &tr)
	if err != nil {
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
	if tr.SettledAt != settledTime {
		t.Errorf("o.SettledAt == %q; want %q", tr.SettledAt, settledTime)
	}
	emittedTime, _ := time.Parse(time.RFC3339, "2018-10-01T00:00:00.000Z")
	if tr.EmittedAt != emittedTime {
		t.Errorf("o.EmittedAt == %q; want %q", tr.EmittedAt, emittedTime)
	}
	if tr.Status != qonto.TransactionStatusCompleted {
		t.Errorf("o.Status == %q; want %q", tr.Status, qonto.TransactionStatusCompleted)
	}
	if tr.Note != "NOTE" {
		t.Errorf("o.Note == %q; want %q", tr.Note, "NOTE")
	}
	if tr.Label != "SOME DEBIT" {
		t.Errorf("o.Label == %q; want %q", tr.Label, "SOME DEBIT")
	}
}
