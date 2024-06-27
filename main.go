package main

import (
	"database/sql"
	"log"

	"github.com/leedrum/simplebank/api"
	db "github.com/leedrum/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(":8080")
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
