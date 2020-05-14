package qonto_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"strings"
	"testing"

	"github.com/ushu/qonto-go/v2"
)

var hasCredentials bool = false
var credentials struct {
	Slug      string
	SecretKey string
}

func init() {
	if b, err := ioutil.ReadFile("./credentials.json"); err == nil {
		if err = json.Unmarshal(b, &credentials); err == nil {
			hasCredentials = true
		}
	}
}

func TestGetOrganization(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	org, err := c.GetOrganization()
	if err != nil {
		t.Fatalf("c.GetOrganization() failed: %v", err)
	}
	if org == nil {
		t.Fatalf("c.GetOrganiation() returned nil Organization")
	}
	if org.Slug != credentials.Slug {
		t.Errorf("org.Slug == %v; want %v", org.Slug, credentials.Slug)
	}
	if len(org.BankAccounts) != 1 {
		t.Errorf("len(org.BankAccounts) == %d; want %d", len(org.BankAccounts), 1)
	}
}

func TestGetBankAccount(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	ba, err := c.GetBankAccount()
	if err != nil {
		t.Fatalf("c.GetBankAccount() failed: %v", err)
	}
	if ba == nil {
		t.Fatalf("c.GetBankAccount() returned nil BankAccount")
	}
	if !strings.HasPrefix(ba.Slug, credentials.Slug) {
		t.Errorf("ba.Slug should start with %s; got %v", credentials.Slug, ba.Slug)
	}
}

func TestGetLabels(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	page, err := c.GetLabels(0, 0)
	if err != nil {
		t.Fatalf("c.GetLabels() failed: %v", err)
	}
	if page == nil {
		t.Fatalf("c.GetLabels() returned nil Labels page")
	}
	if page.Meta.PerPage < len(page.Labels) {
		t.Errorf("want page.Meta.Perpage (%d) >= len(page.Labels) (%d)", page.Meta.PerPage, len(page.Labels))
	}
	if page.Meta.CurrentPage != 1 {
		t.Errorf("page.Meta.CurrentPage == %d; want %d", page.Meta.CurrentPage, 0)
	}
}

func TestGetAllLabels(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	labels, err := c.GetAllLabels(0, 0)
	if err != nil {
		t.Fatalf("c.GetLabels() failed: %v", err)
	}
	if len(labels) == 0 {
		t.Fatalf("c.GetLabels() returned no label")
	}
}

func TestGetMemberships(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	page, err := c.GetMemberships(0, 0)
	if err != nil {
		t.Fatalf("c.GetMemberships() failed: %v", err)
	}
	if page == nil {
		t.Fatalf("c.GetMemberships() returned nil Memberships page")
	}
	if page.Meta.PerPage < len(page.Memberships) {
		t.Errorf("want page.Meta.Perpage (%d) >= len(page.Memberships) (%d)", page.Meta.PerPage, len(page.Memberships))
	}
	if page.Meta.CurrentPage != 1 {
		t.Errorf("page.Meta.CurrentPage == %d; want %d", page.Meta.CurrentPage, 0)
	}
}

func TestGetAllMemberships(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	memberships, err := c.GetAllMemberships(0, 0)
	if err != nil {
		t.Fatalf("c.GetAllMemberships() failed: %v", err)
	}
	if len(memberships) == 0 {
		t.Fatalf("c.GetAllMemberships() returned no membership")
	}
}

func TestGetTransactions(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	ba, err := c.GetBankAccount()
	if err != nil {
		t.Skip()
	}
	page, err := c.GetTransactionsForAccountContext(context.Background(), ba, nil)
	if err != nil {
		t.Fatalf("c.GetTransactionsContext() failed: %v", err)
	}
	if page == nil {
		t.Fatalf("c.GetTransactionsContext() returned nil Memberships page")
	}
}

func TestGetAllTransactions(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	ba, err := c.GetBankAccount()
	if err != nil {
		t.Skip()
	}
	transactions, err := c.GetAllTransactionsForAccount(ba, nil)
	if err != nil {
		t.Fatalf("c.GetAllTransactions() failed: %v", err)
	}
	if len(transactions) == 0 {
		t.Fatalf("c.GetAllTransactions() returned no membership")
	} else {
		t.Logf("Found %d transactions in this account", len(transactions))
	}
}

func TestDownloadAttachment(t *testing.T) {
	t.Parallel()
	if !hasCredentials {
		t.Skip()
	}
	c := qonto.NewClient(credentials.Slug, credentials.SecretKey, nil)
	ba, err := c.GetBankAccount()
	if err != nil {
		t.Skip()
	}
	page, err := c.GetTransactionsForAccount(ba, nil)
	if err != nil {
		t.Skip()
	}

	// lookup an attachment
	transactionWithAttachments := findTransactionWithAttachment(t, page.Transactions)
	if transactionWithAttachments == nil {
		t.Skip("could not find an attachment to download")
	}

	// and finally call the download
	a, err := c.GetAttachment(transactionWithAttachments.AttachmentIDs[0])
	if err != nil {
		t.Fatalf("Could not download attachment: %v", err)
	}
	if exts, _ := mime.ExtensionsByType(a.FileContentType); len(exts) > 0 {
		filename := fmt.Sprintf("./test_attachment%s", exts[0])
		err = c.DownloadAttachmentToFile(transactionWithAttachments.AttachmentIDs[0], filename, os.ModePerm)
		if err != nil {
			t.Errorf("Could not download the attachment: %v", err)
		}
	}
}

func findTransactionWithAttachment(t *testing.T, transactions []*qonto.Transaction) *qonto.Transaction {
	t.Helper()
	for _, t := range transactions {
		if len(t.AttachmentIDs) > 0 {
			return t
		}
	}
	return nil
}
