package main

import (
	"balance_manager/dbmanager"
	"balance_manager/server"
	"database/sql"

	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open("postgres", "host=db port=5432 user=user password=password dbname=test_db sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	db, err := dbmanager.CreateDB(conn)
	if err != nil {
		panic(err)
	}

	s := server.CreateServer(db, ":1337")
	if err := s.Run(); err != nil {
		panic(err)
	}
}
