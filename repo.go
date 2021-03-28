package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

type Repo struct {
	// db *pgx.Conn
	conn *pgx.Conn
}

func NewRepo(conn string) *Repo {
	return &Repo{
		// db: initDB(conn),
		conn: initDB(conn),
	}
}

func initDB(connStr string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// defer conn.Close(context.Background())
	// defer log.Printf("Conn closing")
	err = conn.Ping(context.Background())
	if err != nil {
		log.Printf("Unable to ping database: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Connected to database\n")
	return conn
}
