/*
Package qonto implements an API Client for the Qonto API v2.0.

Basic usage

The package provides a simple Client type for issuing calls to the API.

Example:

   // create a client with valid Credentials
   c := qonto.NewClient("organization-slug", "secret-key")

   // and load the undelying bank account
   ba, _ := c.GetBankAccount()

   // finally we can list all the transactions (the second argument holds optional seach params)
   transactions,_ :=  c.GetAllTransactionsForAccount(ba, nil)

   for _, t := range transactions {
      fmt.Printf("Transaction %s: %f %s", t.ID, t.Amount, t.Currency)
   }
*/
package qonto
