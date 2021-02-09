package main

import (
	"fmt"

	"github.com/roccoblues/form3-account-api"
)

func main() {
	client, _ := form3.NewClient("http://localhost:8080/v1")

	account, _ := client.GetAccount("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	fmt.Println(account.ID)
	fmt.Println(account.CreatedOn)
}
