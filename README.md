# Qonto API Client for go

[![Build Status](https://travis-ci.org/ushu/qonto-go.svg?branch=master)](https://travis-ci.org/ushu/qonto-go)

This library is intended to connect to the [Qonto API v2.0] to access the details of the user's transactions.

*Note*: there is another [go package](https://github.com/toorop/go-qonto) that you might want to consider for consuming this API. I designed this variant to be able to easily use `context.Context` for parallel API calls and use a custom `http.Client` for AppEngine, so it will depend on your use case whether or not it is relevant.

## Getting Started

```sh
go get github.com/ushu/qonto-go
```

### Usage

The simple `Client` type allows to connect and fetch information from the API:

```go
 // create a client with valid Credentials
 c := qonto.NewClient("organization-slug", "secret-key")

 // load organization details
 o, _ := c.GetOrganization()

 // and read bank account details
 a := o.BankAccounts[0]
 t, _, _ := c.GetTransactions(a.Slug, a.IBAN, nil)
```

For more advanced use cases a `ContextClient` type is provided.
Compared to the `Client` type, it handles two more advanced use cases:

- providing a custom `http.Client`, which will be useful on AppEngine for eg
- connecting to several accounts using the same clients: the credentials are passed through the provided context, so you can use as many credentils as necessry, in parallel, using the same client

Example:

```go
// we suppose we are on AppEngine
ctx := appengine.NewContext(r)
// and provide a custom http.Client
httpClient := urlfetch.Client(ctx)

// then we obtain a custom ContextClient
c := &qonto.ContextClient(httpClient)

// to call the API, we must embed our credentials into the context
ctx = qonto.WithCredentials(ctx, "organization-slug", "secret-key")

// and pass the context with credentials to the calls
o, _ := c.GetOrganization(ctx)
```

## Contributing

Feel free to contribute anytime !

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

[Qonto API v2.0]: https://api-doc.qonto.eu/2.0
[Stripe API]: https://api-doc.qonto.eu/2.0
