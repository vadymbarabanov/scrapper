package main

import (
	"fmt"

	"github.com/vadymbarabanov/scrapper/account"
)

func main() {
	newAcc := account.CreateAccount("lynn")
	fmt.Println(newAcc)

	newAcc.Deposit(30)

	err := newAcc.Withdraw(10)
	if err != nil {
		fmt.Println(err)
	}

	newAcc.ChangeOwner("bob")
	fmt.Println(newAcc)
}
