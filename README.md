# [Aspire Mini API](https://github.com/avinragh/aspire)

API in Go to add and view loans and payments

## Author
Avinash Raghunathan

## Features

* Implements Create, Fetch and Delete APIs
* Suitable for use in other software projects
* Tests that run from `docker-compose up`
* Simple and Concise
* Constants available for ISO standard format values e.g: Country, Currency and Bank Codes
* User Defined context and Http client can be passed to have more control over the API Operations like Timeout and Cancellation

## Usage

### Create Account

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/avinragh/form3-sdk/account"
	"github.com/avinragh/form3-sdk/form3"
	"github.com/davecgh/go-spew/spew" //for pretty printing struct values
)

func main() {
	accountData := &account.AccountData{
		Type:           "accounts",
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Attributes: account.NewAccountAttributes().
			WithCountry(form3.CountryUnitedKingdom).
			WithBic("NWBKGB22").
			WithBankID("400300").
			WithBankIDCode(form3.BankIDCodeUnitedKingdom).
			WithAccountNumber("41426815").
			WithBaseCurrency(form3.CurrencyUnitedKingdom).
			WithJointAccount(false).
			WithName("Avinash Raghunathan").
			WithAlternativeNames(
				"Raghunathan", "Avi",
			).
			WithAccountMatchingOptOut(true),
	}

	accountInstance := account.NewAccount(accountData)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := account.NewClient(&http.Client{}, form3.V1)
	createdAccount, err := client.CreateAccount(ctx, accountInstance)
	if err != nil {
		fmt.Println(err)
		if createdAccount != nil {
			spew.Dump(createdAccount.Error)
		}
	}
	spew.Dump(createdAccount)
```

### Fetch Account

```go
	fetchedAccount, err := client.FetchAccount(ctx, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	if err != nil {
		fmt.Println(err)
		if fetchedAccount != nil {
			spew.Dump(fetchedAccount.Error)
		}
	}
	spew.Dump(fetchedAccount)
```

### Delete Account

```go
	var version int64
	version = 0
    	err:= client.DeleteAccount(ctx,"ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",version)
    	if err != nil {
        	fmt.Println(err)
    	}
```

### Run Tests

```bash
docker-compose up
```
or Run in the Background

```bash
docker-compose up -d
```
```bash
docker-compose logs accountapi-sdk
```



