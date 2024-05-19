package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAllAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		phone_number VARCHAR(15) NOT NULL,
		balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	return err
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=root dbname=postgres password=12345 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO accounts (
		first_name,
		last_name,
		phone_number,
		balance,
		created_at
	) VALUES (
		$1, $2, $3, $4, $5
	);`

	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.PhoneNumber, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = $1`

	_, err := s.db.Query(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	var rows *sql.Rows
	var err error

	query := `SELECT * FROM accounts WHERE id = $1`

	if rows, err = s.db.Query(query, id); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return ScanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account id %v not found", id)
}

func (s *PostgresStore) GetAllAccounts() ([]*Account, error) {
	var rows *sql.Rows
	var err error
	var accounts []*Account

	if rows, err = s.db.Query(`SELECT * FROM accounts`); err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account *Account

		if account, err = ScanIntoAccount(rows); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return accounts, nil
}

func ScanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	if err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.PhoneNumber,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return account, nil
}
