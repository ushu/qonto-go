package main

import (
	"github.com/ushu/qonto-go"
	"fmt"
	"encoding/json"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: test SLUG SECRET_KEY")
		os.Exit(1)
	}
	c := qonto.NewClient(os.Args[1], os.Args[2])

	// Load organization
	o, _ := c.GetOrganization()
	j, _ := json.MarshalIndent(o, "", "  ")

	fmt.Println("Organization")
	fmt.Println("============")
	fmt.Println()
	fmt.Println(string(j))
	fmt.Println()

	a := o.BankAccounts[0]
	t, _, _ := c.GetTransactions(a.Slug, a.IBAN, nil)
	j, _ = json.MarshalIndent(t, "", "  ")

	fmt.Println("Transactions")
	fmt.Println("============")
	fmt.Println()
	fmt.Println(string(j))
	fmt.Println()
}
