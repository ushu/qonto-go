package qonto

import (
	"context"
	"errors"
)

type Credentials struct {
	Slug      string // the unique slug for the organization
	SecretKey string // the secret key to Authenticate to Qonto
}

var DefaultCredentials *Credentials

const credentialsContextKey = "github.com/ushu/qonto-go.Credentials"

// MissingCredentialsError is returned when trying to call the client without
// adding credentials into the parent context.
var MissingCredentialsError = errors.New("Missing Qonto credentials")

func WithCredentials(ctx context.Context, slug, secretKey string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, credentialsContextKey, &Credentials{
		Slug:      slug,
		SecretKey: secretKey,
	})
}

func GetCredentials(ctx context.Context) (string, string, error) {
	c := ctx.Value(credentialsContextKey)
	if c == nil {
		c = DefaultCredentials
	}
	if c == nil {
		return "", "", MissingCredentialsError
	}
	cred := c.(*Credentials)
	return cred.Slug, cred.SecretKey, nil
}

func SetDefaultCredentials(slug, secretKey string) {
	DefaultCredentials = &Credentials{
		Slug:      slug,
		SecretKey: secretKey,
	}
}
