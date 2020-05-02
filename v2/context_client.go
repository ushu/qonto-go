package qonto

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// BackendURL holds the base URL for calling Qonto API v2.
const BackendURL = "https://thirdparty.qonto.eu/v2/"

// A ContextClient is an API client for the Qonto API v2.
// It can be passed a http.Client instance, but its zero value uses http.DefaultClient.
type ContextClient struct {
	HTTPClient *http.Client
}

// GetOrganization fetches the organization details for the provided credentials.
// It reads its credentials from ctx, meaning it expects WithCredentials to be called before to load
// them into the context.
func (c *ContextClient) GetOrganization(ctx context.Context) (*Organization, error) {
	// we need to build a path with the organization slug: ".../organizations/my-slug-here"
	slug, _, err := GetCredentials(ctx)
	if err != nil {
		return nil, err
	}
	path := "organizations/" + url.PathEscape(slug)

	// call Qonto
	buf, err := c.authenticatedGet(ctx, path)
	if err != nil {
		return nil, err
	}

	// and decode the Organization (JSON)
	var o OrganizationResponse
	err = json.Unmarshal(buf, &o)
	return o.Organization, err
}

// GetTransactions fetches the transaction details for the provided slug and iban.
// It reads its credentials from ctx, meaning it expects WithCredentials to be called before to load
// them into the context.
// The optional opt parameter allows to define pagination options.
func (c *ContextClient) GetTransactions(ctx context.Context, slug, iban string, opt *PaginationOptions) ([]*Transaction, *ResponseMeta, error) {
	// we need to path both the slug and IBAN as query params
	v := url.Values{}
	v.Set("slug", slug)
	v.Set("iban", iban)
	if opt != nil {
		if opt.CurrentPage > 0 {
			v.Set("current_page", strconv.Itoa(opt.CurrentPage))
		}
		if opt.PerPage > 0 {
			v.Set("per_page", strconv.Itoa(opt.PerPage))
		}
	}
	path := "transactions?" + v.Encode()

	// call Qonto
	buf, err := c.authenticatedGet(ctx, path)
	if err != nil {
		return nil, nil, err
	}

	// and decode the Transactions (JSON)
	var o TransactionsResponse
	err = json.Unmarshal(buf, &o)
	return o.Transactions, o.Meta, err
}

// IterTransactionPages creates a new ContextIterator for the provided iban and slug.
func (c *ContextClient) IterTransactionPages(ctx context.Context, slug, iban string, perPage int) *ContextIterator {
	return &ContextIterator{
		c:        c,
		Slug:     slug,
		IBAN:     iban,
		nextPage: 0,
		perPage:  perPage,
		done:     false,
	}
}

func (c *ContextClient) authenticatedGet(ctx context.Context, path string) ([]byte, error) {
	slug, secretKey, err := GetCredentials(ctx)
	if err != nil {
		return nil, err
	}

	// prepare the Get Request
	u := BackendURL
	if path != "" {
		u = u + "/" + path
	}
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.Header.Set("Authorization", slug+":"+secretKey)
	req = req.WithContext(ctx)

	// call the endpoint
	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// and dump all data from the response
	buf, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return buf, err
}
