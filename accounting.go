package nicehash

import (
	"fmt"

	"github.com/mineselskabet/go-bitcoin"
)

type Currency string

type Activity string

const (
	DEPOSIT       Activity = "DEPOSIT"
	WITHDRAWAL    Activity = "WITHDRAWAL"
	HASHPOWER     Activity = "HASHPOWER"
	MINING        Activity = "MINING"
	EXCHANGE      Activity = "EXCHANGE"
	UNPAID_MINING Activity = "UNPAID_MINING"
	OTHER         Activity = "OTHER"
)

type ActivityCompletion string

const (
	COMPLETED ActivityCompletion = "COMPLETED"
	OPEN      ActivityCompletion = "OPEN"
	ALL       ActivityCompletion = "ALL"
)

type Transaction struct {
	Id       string         `json:"id"`
	Amount   bitcoin.Amount `json:"amount"`
	Fee      bitcoin.Amount `json:"feeAmount"`
	Time     Time           `json:"time"`
	Type     Activity       `json:"type"`
	Currency Currency       `json:"activityCurrency"`
}

type Balance struct {
	Active       bool
	TotalBalance bitcoin.Amount
	Available    bitcoin.Amount
	Debt         bitcoin.Amount
	Pending      bitcoin.Amount
}

func (c *Client) AccountingBalance(currency Currency) (*Balance, error) {
	URL := fmt.Sprintf("https://api2.nicehash.com/main/api/v2/accounting/account2/%s", currency)

	var balance Balance

	err := c.exchangeJSON("GET", URL, nil, &balance)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (c *Client) AccountingActivity(currency Currency, typ Activity, stage ActivityCompletion) ([]Transaction, error) {
	URL := fmt.Sprintf("https://api2.nicehash.com/main/api/v2/accounting/activity/%s?stage=%s&type=%s", currency, stage, typ)

	var list []Transaction

	err := c.exchangeJSON("GET", URL, nil, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
