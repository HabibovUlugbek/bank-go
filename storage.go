package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdatAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByPhoneNumber(int64) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres port=1000 password=gobank sslmode=disable"
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

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		password varchar(50),
		balance serial, 
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
	insert into account (first_name, last_name, number, password, balance, created_at) 
	values ($1, $2, $3, $4, $5, $6)`

	res, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Password, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `delete from account where id = $1`

	_, err := s.db.Query(query, id)
	return err
}

func (s *PostgresStore) UpdatAccount(acc *Account) error {
	query := `update account set first_name = $1, last_name = $2, number = $3, balance = $4, where id = $5`

	_, err := s.db.Query(query, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.ID)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `select id, first_name, last_name, number, balance, created_at from account`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
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

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `select id, first_name, last_name, number, balance, created_at from account where id = $1`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccountByPhoneNumber(phoneNumber int64) (*Account, error) {
	query := `select id, first_name, last_name, number, balance, created_at from account where number = $1`

	rows, err := s.db.Query(query, phoneNumber)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", phoneNumber)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}
