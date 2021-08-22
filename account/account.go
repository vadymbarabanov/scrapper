package account

import (
	"errors"
	"fmt"
)

var errNoMoney = errors.New("Error: Cannot withdraw, not enough money.")

type Account struct {
	owner   string
	balance int
}

func (a *Account) Deposit(amount int) {
	a.balance += amount
}

func (a Account) GetBalance() int {
	return a.balance
}

func (a Account) GetOwner() string {
	return a.owner
}

func (a Account) String() string {
	return fmt.Sprintf("Owner: %s\nBalance: %d\n", a.GetOwner(), a.GetBalance())
}

func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

func CreateAccount(name string) *Account {
	acc := Account{owner: name, balance: 0}
	return &acc
}
