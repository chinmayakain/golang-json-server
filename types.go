package main

import (
	"math/rand"
	"time"
)

type TransferReq struct {
	ToAccount string `json:"toAccount"`
	Amount string `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber int64     `json:"phoneNumber"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: rand.Int63n(10000000000),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}
