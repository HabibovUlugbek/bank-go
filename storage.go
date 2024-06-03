package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdatAccount(*Account) error
	GetAccountById(int) (*Account, error)
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
		balance serial, 
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account (id, first_name, last_name, number, balance, created_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(query, acc.ID, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `delete from account where id = $1`

	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) UpdatAccount(acc *Account) error {
	query := `update account set first_name = $1, last_name = $2, number = $3, balance = $4, where id = $5`

	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.ID)
	return err
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `select id, first_name, last_name, number, balance, created_at from account where id = $1`

	row := s.db.QueryRow(query, id)

	acc := &Account{}

	err := row.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &Account{}, nil
}
