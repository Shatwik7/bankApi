package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	Login(int, int64) (*Account, error)
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(int) (bool, error)
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Login(int, int64) (*Account, error) {
	return nil, nil
}
func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}
func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=admin dbname=bank password=password sslmode=disable"
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
	encrypted_password VARCHAR(255) NOT NULL,
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
	INSERT INTO account (first_name, last_name, encrypted_password, number, balance,created_at) 
	VALUES ($1, $2, $3, $4,$5,$6) RETURNING id;
	`
	var id int
	err := s.db.QueryRow(query, account.FirstName, account.LastName, account.EncryptedPassword, account.Number, account.Balance, account.CreatedAt).Scan(&id)
	if err != nil {
		return &Account{}, fmt.Errorf("error inserting account: %w", err)
	}

	account.ID = id // Set the generated ID
	return account, nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

// func (s *PostgresStore) TransferAmount(Account_id int, amount int) (bool, error) {
// 	query:=``
// }

func (s *PostgresStore) DeleteAccount(id int) (bool, error) {
	query := `DELETE FROM account WHERE id=$1 ;`
	_, err := s.db.Query(query, id)
	if err != nil {
		return false, err
	}
	return true, err
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `SELECT * FROM account WHERE id=$1;`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("ACCOUNT NOT FOUND")
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`
	rows, err := s.db.Query(query)
	if err != nil {
		return []*Account{}, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	if err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.EncryptedPassword, // scaning is relative to the sql schema
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	); err != nil {
		return nil, err
	}
	return account, nil
}
