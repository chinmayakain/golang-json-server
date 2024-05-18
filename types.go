package main

import "math/rand"

type Account struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber int64  `json:"phoneNumber"`
	Balance     int64  `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:          rand.Int63n(10000),
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: rand.Int63n(10000000000),
	}
}
