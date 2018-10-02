package qonto

// Client allows to send requests to the Qonto API servers.
type Client struct {
	c         *ContextClient
	Slug      string // the organization slug.
	SecretKey string // the secret key, associated to the organization.
}

// NewClient creates and initialisez a new Client with the provided credentials.
func NewClient(slug, secretKey string) *Client {
	return &Client{
		c:         &ContextClient{},
		Slug:      slug,
		SecretKey: secretKey,
	}
}

// GetOrganization fetches the organization details.
func (c *Client) GetOrganization() (*Organization, error) {
	ctx := WithCredentials(nil, c.Slug, c.SecretKey)
	return c.c.GetOrganization(ctx)
}

// GetTransactions fetches the transaction details for the provided slug and iban.
// The optional opt parameter allows to define pagination options.
func (c *Client) GetTransactions(slug, iban string, opt *PaginationOptions) ([]*Transaction, *ResponseMeta, error) {
	ctx := WithCredentials(nil, c.Slug, c.SecretKey)
	return c.c.GetTransactions(ctx, slug, iban, opt)
}

// IterTransactionPages creates a new Iterator for the provided iban and slug.
func (c *Client) IterTransactionPages(slug, iban string, perPage int) *Iterator {
	ctx := WithCredentials(nil, c.Slug, c.SecretKey)
	i := c.c.IterTransactionPages(ctx, slug, iban, perPage)
	return &Iterator{
		c: c,
		i: i,
	}
}
