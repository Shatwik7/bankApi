package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}
func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    number BIGINT UNIQUE NOT NULL,
    balance BIGINT NOT NULL,
    created_at timestamp);`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateAccount(account *Account) (*Account, error) {
	query := `
	INSERT INTO account (first_name, last_name, number, balance,created_at) 
	VALUES ($1, $2, $3, $4,$5) RETURNING id;
	`
	var id int
	err := s.db.QueryRow(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt).Scan(&id)
	if err != nil {
		return &Account{}, fmt.Errorf("error inserting account: %w", err)
	}

	account.ID = id // Set the generated ID
	return account, nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`
	rows, err := s.db.Query(query)
	if err != nil {
		return []*Account{}, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		if err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
