package qonto

import (
	"context"
	"errors"
)

// ErrDone is returned by iterator if Next() is called after the last page.
var ErrDone = errors.New("no more page")

// The PaginationOptions holds the pagination details for the requests (pagination info in
// response if held by ResponseMeta).
type PaginationOptions struct {
	CurrentPage int `json:"current_page"` // page to start listing
	PerPage     int `json:"per_page"`     // maximum number of entries returned
}

// ContextIterator is an iterator type for the ContextClient.
type ContextIterator struct {
	c        *ContextClient
	Slug     string // the slug of the bank account
	IBAN     string // the iban of the bank account
	nextPage int    // the page that will be fetch on the next call to Next()
	perPage  int    // the requested number of transactions per page
	done     bool   // did we load all the pages ?
}

// Next fetches and returns the next page of transactions.
func (i *ContextIterator) Next(ctx context.Context) ([]*Transaction, bool, error) {
	if i.done {
		return nil, true, ErrDone
	}

	opt := &PaginationOptions{
		CurrentPage: i.nextPage,
		PerPage:     i.perPage,
	}
	t, m, err := i.c.GetTransactions(ctx, i.Slug, i.IBAN, opt)
	if err != nil {
		return nil, false, err
	}

	if m.NextPage != nil {
		i.nextPage = *m.NextPage
	} else {
		i.done = true
	}
	return t, i.done, nil
}

// Iterator allows to iterate pages of Transactions returned by the Client.
type Iterator struct {
	c *Client
	i *ContextIterator
}

// Next fetches and returns the next page of transactions.
func (i *Iterator) Next() ([]*Transaction, bool, error) {
	ctx := WithCredentials(nil, i.c.Slug, i.c.SecretKey)
	return i.i.Next(ctx)
}
