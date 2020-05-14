package qonto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// BaseURL is the root URL for all API calls
var BaseURL = "https://thirdparty.qonto.eu/v2"

// ErrMissingBankAccountSlug error
var ErrMissingBankAccountSlug = errors.New("Missing \"slug\" parameter")

// ErrMissingBankAccountIBAN error
var ErrMissingBankAccountIBAN = errors.New("Missing \"iban\" parameter")

// ErrBankAccountNeeded error
var ErrBankAccountNeeded = errors.New("Cannot pass a nil bank account")

// ErrAttachementNeeded error
var ErrAttachementNeeded = errors.New("Cannot pass a nil attachment")

// ErrMissingAttachmentURL error
var ErrMissingAttachmentURL = errors.New("This attachment as no download URL")

// APIError holds the error information sent in API responses.
type APIError struct {
	// Request the outgoing API Request
	Request *http.Request `json:"-"`
	// Response the incoming API Response
	Response *http.Response `jsoon:"-"`
	// Message the error description returned by Qonto
	Message string `json:"message"`
}

func (q APIError) Error() string {
	return fmt.Sprintf("%s %s ➡︎ %1d \"%s\" ", q.Request.Method, q.Request.URL.String(), q.Response.StatusCode, q.Message)
}

// Client allows to send requests to the Qonto API servers.
type Client struct {
	h         *http.Client
	Slug      string // the organization slug.
	SecretKey string // the secret key, associated to the organization.
}

// NewClient creates and initialisez a new Client with the provided credentials.
func NewClient(slug, secretKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		h:         httpClient,
		Slug:      slug,
		SecretKey: secretKey,
	}
}

// GetOrganization fetches the organization details.
//
// The API docs (https://api-doc.qonto.eu/2.0/organizations/show-organization-1) mention an {id} param
// but as of today Qonto supports only one organization per account, which id equals the
// authentication slug.
func (c *Client) GetOrganization() (*Organization, error) {
	return c.GetOrganizationContext(context.Background())
}

// GetOrganizationContext fetches the organization details, attaching ctx to the request.
func (c *Client) GetOrganizationContext(ctx context.Context) (*Organization, error) {
	path := fmt.Sprintf("%s/organizations/%s", BaseURL, c.Slug)

	// this endpoint responds with a JSON object holding an "organization" key
	var response struct {
		Organization Organization `json:"organization"`
	}
	if err := c.getJSON(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response.Organization, nil
}

// GetBankAccount returns the bank account for the Organization.
//
// As of today, Qonto only supports one bank account per Organization.
func (c *Client) GetBankAccount() (*BankAccount, error) {
	return c.GetBankAccountContext(context.Background())
}

// GetBankAccountContext returns the bank account for the Organization.
func (c *Client) GetBankAccountContext(ctx context.Context) (*BankAccount, error) {
	org, err := c.GetOrganizationContext(ctx)
	if err != nil {
		return nil, err
	}
	if len(org.BankAccounts) == 0 {
		return nil, fmt.Errorf("Found no bank account for Organization %s", c.Slug)
	}
	return org.BankAccounts[0], nil
}

// LabelsPage represents the data returned by Qonto when calling GetLabels*
type LabelsPage struct {
	Labels []Label `json:"labels"`
	Meta   Meta    `json:"meta"`
}

// GetLabels fetches the list of labels defined in the current Organization
func (c *Client) GetLabels(currentPage, perPage int) (page *LabelsPage, err error) {
	return c.GetLabelsContext(context.Background(), currentPage, perPage)
}

// GetLabelsContext fetches the list of labels defined in the current Organization
func (c *Client) GetLabelsContext(ctx context.Context, currentPage, perPage int) (page *LabelsPage, err error) {
	u, err := addPaginationQueryParams(BaseURL+"/labels", currentPage, perPage)
	if err != nil {
		return nil, err // ⬅︎ should not happen unless BaseURL is modified
	}
	err = c.getJSON(ctx, u, &page)
	return
}

// GetAllLabels fetches the list of labels defined in the current Organization
func (c *Client) GetAllLabels(startPage, perPage int) (labels []Label, err error) {
	return c.GetAllLabelsContext(context.Background(), startPage, perPage)
}

// GetAllLabelsContext fetches the list of labels defined in the current Organization
func (c *Client) GetAllLabelsContext(ctx context.Context, startPage, perPage int) (labels []Label, err error) {
	var res *LabelsPage
	for {
		res, err = c.GetLabelsContext(ctx, startPage, perPage)
		if err != nil {
			break
		}
		labels = append(labels, res.Labels...)
		if res.Meta.NextPage == nil {
			break
		} else {
			startPage = *res.Meta.NextPage
		}
	}
	return
}

// MembershipsPage represents the data returned by Qonto when calling GetMemberships*
type MembershipsPage struct {
	Memberships []Membership `json:"memberships"`
	Meta        Meta         `json:"meta"`
}

// GetMemberships fetches the list of members of the current Organization
func (c *Client) GetMemberships(currentPage, perPage int) (page *MembershipsPage, err error) {
	return c.GetMembershipsContext(context.Background(), currentPage, perPage)
}

// GetMembershipsContext fetches the list of members of the current Organization
func (c *Client) GetMembershipsContext(ctx context.Context, currentPage, perPage int) (page *MembershipsPage, err error) {
	u, err := addPaginationQueryParams(BaseURL+"/memberships", currentPage, perPage)
	if err != nil {
		return nil, err // ⬅︎ should not happen unless BaseURL is modified
	}
	err = c.getJSON(ctx, u, &page)
	return
}

// GetAllMemberships fetches the list of members of the current Organization
func (c *Client) GetAllMemberships(startPage, perPage int) (memberships []Membership, err error) {
	return c.GetAllMembershipsContext(context.Background(), startPage, perPage)
}

// GetAllMembershipsContext fetches the list of members of the current Organization
func (c *Client) GetAllMembershipsContext(ctx context.Context, startPage, perPage int) (memberships []Membership, err error) {
	var res *MembershipsPage
	for {
		res, err = c.GetMembershipsContext(ctx, startPage, perPage)
		if err != nil {
			break
		}
		memberships = append(memberships, res.Memberships...)
		if res.Meta.NextPage == nil {
			break
		} else {
			startPage = *res.Meta.NextPage
		}
	}
	return
}

// TransactionsPage represents the data returned by Qonto when calling GetTransactions*
type TransactionsPage struct {
	Transactions []*Transaction `json:"transactions"`
	Meta         Meta           `json:"meta"`
}

// GetTransactionOptions list all the (optional) search options for the Get[All]Transactions* calls
type GetTransactionOptions struct {
	Statuses      []TransactionStatus
	UpdatedAtFrom *time.Time
	UpdatedAtTo   *time.Time
	SettledAtFrom *time.Time
	SettledAtTo   *time.Time
	SortBy        *string
	CurrentPage   *int
	PerPage       *int
}

// GetTransactionsForAccount fetches the list of transactions of the given bank account
func (c *Client) GetTransactionsForAccount(ba *BankAccount, options *GetTransactionOptions) (page *TransactionsPage, err error) {
	if ba == nil {
		return nil, ErrBankAccountNeeded
	}
	return c.GetTransactions(ba.Slug, ba.IBAN, options)
}

// GetTransactions fetches the list of transactions of the given bank account
func (c *Client) GetTransactions(bankAccountID, IBAN string, options *GetTransactionOptions) (page *TransactionsPage, err error) {
	return c.GetTransactionsContext(context.Background(), bankAccountID, IBAN, options)
}

// GetTransactionsForAccountContext fetches the list of transactions of the given bank account
func (c *Client) GetTransactionsForAccountContext(ctx context.Context, ba *BankAccount, options *GetTransactionOptions) (page *TransactionsPage, err error) {
	if ba == nil {
		return nil, ErrBankAccountNeeded
	}
	return c.GetTransactionsContext(ctx, ba.Slug, ba.IBAN, options)
}

// GetTransactionsContext fetches the list of transactions of the given bank account
func (c *Client) GetTransactionsContext(ctx context.Context, bankAccountID, IBAN string, options *GetTransactionOptions) (page *TransactionsPage, err error) {
	u, err := getTransactionsURL(bankAccountID, IBAN, options)
	if err != nil {
		return nil, err // ⬅︎ should not happen unless BaseURL is modified
	}
	err = c.getJSON(ctx, u, &page)
	return
}

// GetAllTransactionsForAccount fetches the list of transactions of the given bank account
func (c *Client) GetAllTransactionsForAccount(ba *BankAccount, options *GetTransactionOptions) (transactions []*Transaction, err error) {
	if ba == nil {
		return nil, ErrBankAccountNeeded
	}
	return c.GetAllTransactions(ba.Slug, ba.IBAN, options)
}

// GetAllTransactions fetches the list of transactions of the given bank account
func (c *Client) GetAllTransactions(bankAccountID, IBAN string, options *GetTransactionOptions) (transactions []*Transaction, err error) {
	return c.GetAllTransactionsContext(context.Background(), bankAccountID, IBAN, options)
}

// GetAllTransactionsForAccountContext fetches the list of transactions of the given bank account
func (c *Client) GetAllTransactionsForAccountContext(ctx context.Context, ba *BankAccount, options *GetTransactionOptions) (transactions []*Transaction, err error) {
	if ba == nil {
		return nil, ErrBankAccountNeeded
	}
	return c.GetAllTransactionsContext(ctx, ba.Slug, ba.IBAN, options)
}

// GetAllTransactionsContext fetches the list of transactions of the given bank account
func (c *Client) GetAllTransactionsContext(ctx context.Context, bankAccountID, IBAN string, options *GetTransactionOptions) (transactions []*Transaction, err error) {
	var res *TransactionsPage

	// we keep track of the current page and increment it one by one
	var currentPage = 1
	if options != nil && options.CurrentPage != nil {
		currentPage = *options.CurrentPage
	}

	for {
		var callOptions GetTransactionOptions
		if options != nil {
			callOptions = *options // ⬅︎ we copy all existing options
		}
		callOptions.CurrentPage = &currentPage // ⬅︎ and overwrite the current Page

		res, err = c.GetTransactionsContext(ctx, bankAccountID, IBAN, &callOptions)
		if err != nil {
			break
		}
		transactions = append(transactions, res.Transactions...)
		if res.Meta.NextPage == nil {
			break
		} else {
			currentPage = *res.Meta.NextPage
		}
	}
	return
}

// GetAttachment downloads a remote attachment given it's id
func (c *Client) GetAttachment(id string) (attachment *Attachment, err error) {
	return c.GetAttachmentContext(context.Background(), id)
}

// GetAttachmentContext downloads a remote attachment given it's id
func (c *Client) GetAttachmentContext(ctx context.Context, id string) (*Attachment, error) {
	u := fmt.Sprintf("%s/attachments/%s", BaseURL, id)

	var response struct {
		Attachment *Attachment `json:"attachment"`
	}
	err := c.getJSON(ctx, u, &response)
	if err != nil {
		return nil, err
	}
	return response.Attachment, nil
}

// DownloadAttachment downloads the file contents of an attachment
func (c *Client) DownloadAttachment(a *Attachment) ([]byte, error) {
	return c.DownloadAttachmentContext(context.Background(), a)
}

// DownloadAttachmentContext downloads the file contents of an attachment
func (c *Client) DownloadAttachmentContext(ctx context.Context, a *Attachment) ([]byte, error) {
	if a == nil {
		return nil, ErrAttachementNeeded
	}
	if a.URL == "" {
		return nil, ErrMissingAttachmentURL
	}
	req, err := http.NewRequestWithContext(ctx, "GET", a.URL, nil)
	res, err := c.h.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		_ = res.Body.Close()
		return nil, err
	}
	return data, res.Body.Close()
}

// DownloadAttachmentToFile downloads the file contents of an attachment into a local file
func (c *Client) DownloadAttachmentToFile(id, filename string, perm os.FileMode) error {
	return c.DownloadAttachmentToFileContext(context.Background(), id, filename, perm)
}

// DownloadAttachmentToFileContext downloads the file contents of an attachment into a local file
func (c *Client) DownloadAttachmentToFileContext(ctx context.Context, id, filename string, perm os.FileMode) error {
	a, err := c.GetAttachmentContext(ctx, id)
	if err != nil {
		return err
	}
	data, err := c.DownloadAttachmentContext(ctx, a)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, perm)
}

func addPaginationQueryParams(baseURL string, currentPage, perPage int) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL, err // ⬅︎ should not happen unless BaseURL is modified
	}

	// encode the options in a query params
	query := url.Values{}
	if currentPage > 0 {
		query.Set("current_page", strconv.Itoa(currentPage))
	}
	if perPage > 0 {
		query.Set("per_page", strconv.Itoa(perPage))
	}
	u.RawQuery = u.Query().Encode()

	return u.String(), nil
}

func getTransactionsURL(slug, IBAN string, options *GetTransactionOptions) (string, error) {
	u, err := url.Parse(BaseURL + "/transactions")
	if err != nil {
		return "", err // ⬅︎ should not happen unless BaseURL is modified
	}

	// from here we append lots of query params
	query := url.Values{}

	// mandatory params
	if slug == "" {
		return "", ErrMissingBankAccountSlug
	}
	query.Set("slug", slug)
	if IBAN == "" {
		return "", ErrMissingBankAccountIBAN
	}
	query.Set("iban", IBAN)

	// optional params
	if options != nil {
		for _, s := range options.Statuses {
			query.Add("status", string(s))
		}
		if options.UpdatedAtFrom != nil {
			query.Set("updated_at_from", options.UpdatedAtFrom.Format(time.RFC3339))
		}
		if options.UpdatedAtTo != nil {
			query.Set("updated_at_to", options.UpdatedAtTo.Format(time.RFC3339))
		}
		if options.SettledAtFrom != nil {
			query.Set("settled_at_from", options.SettledAtFrom.Format(time.RFC3339))
		}
		if options.SettledAtTo != nil {
			query.Set("settled_at_to", options.SettledAtTo.Format(time.RFC3339))
		}
		if options.SortBy != nil {
			query.Set("sort_by", *options.SortBy)
		}
		if options.CurrentPage != nil && *options.CurrentPage > 0 {
			query.Set("current_page", strconv.Itoa(*options.CurrentPage))
		}
		if options.PerPage != nil && *options.PerPage > 0 {
			query.Set("per_page", strconv.Itoa(*options.PerPage))
		}
	}

	// finally we encode all the query params in the URL
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (c *Client) getJSON(ctx context.Context, u string, ref interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return err // ⬅︎ should not happen unless we override BaseURL
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s:%s", c.Slug, c.SecretKey))
	res, err := c.h.Do(req)
	if err != nil {
		return fmt.Errorf("Qonto API could not be reached: %w", err)
	}
	if res.StatusCode > 299 {
		ae := APIError{
			Request:  req,
			Response: res,
		}
		err = json.NewDecoder(res.Body).Decode(&ae)
		_ = res.Body.Close()
		if err == nil {
			return ae // ⬅︎ APIError will retain the description sent by Qonto
		}
		// could not decode the JSON body, we send a generic error
		return fmt.Errorf("GET %s returned %d", u, res.StatusCode)
	}
	err = json.NewDecoder(res.Body).Decode(ref)
	if err != nil {
		_ = res.Body.Close()
		return fmt.Errorf("Could not decode the response from Qonto API: %w", err)
	}
	return res.Body.Close()
}
