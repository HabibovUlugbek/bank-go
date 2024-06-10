package main

import (
	"log"
)

func seedAccount(store Storage, fname, lname, pw string) {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}
}

func seedAccounts(s Storage) {
	seedAccount(s, "johnatahan", "bale", "inesta")
}

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	seedAccounts(store)
	server := NewAPIServer(":3000", store)
	server.Run()
}
