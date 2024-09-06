package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mikeheft/go-backend/api"
	db "github.com/mikeheft/go-backend/db/sqlc"
	"github.com/mikeheft/go-backend/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannont connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

}
