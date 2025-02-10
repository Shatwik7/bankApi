package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Token     string `json:"token"`
	Balance   int64  `json:"balance"`
}
type TrasnferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}
type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}
type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"encrypted_password"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string, password string) (*Account, error) {
	Enpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:                rand.Intn(10000),
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(10000)),
		Balance:           1000,
		EncryptedPassword: string(Enpass),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
