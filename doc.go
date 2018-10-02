/* 
Package qonto implements an API Client for the Qonto API v2.0.

Basic usage

The package provides a simple Client type for issuing calls to the API.

Example:

   // create a client with valid Credentials
   c := qonto.NewClient("organization-slug", "secret-key")

   // load organization details
   o, _ := c.GetOrganization()

   // and read bank account details
   a := o.BankAccounts[0]
   t, _, _ := c.GetTransactions(a.Slug, a.IBAN, nil)

Advanced usage

The package also offers a more advanced version of the Client named ContextClient 
that uses context.Context on each call and allows to use a custom http.Client, for
example when deploying the code on AppEngine.

Note that in this case the credentials must be enbedded into the provided Context, using
the WithCredentials function.

Example:

   // in this exemple we use an AppEngine http.Client
   ctx := appengine.NewContext(r)
   httpClient := urlfetch.Client(ctx)

   // to obtain a custom ContextClient
   c := &qonto.ContextClient(httpClient)

   // to call the API, we must embed our credentials into the context
   ctx = qonto.WithCredentials(ctx, "organization-slug", "secret-key")

   // and pass the context with credentials to the calls
   o, _ := c.GetOrganization(ctx)

This approach allows the ContextClient to handle several scenarios: use of custom HTTP 
clients libraries (as AppEngine above), connection to multiple accounts using a single client, 
handling of parallel execution/cancelation with the use of Context etc.
*/
package qonto
