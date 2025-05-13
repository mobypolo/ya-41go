package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// явно несуществующий домен
	dsn := "postgres://user:password@doesnotexist.zzz:5432/dbname?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("sql.Open failed: %v", err)
	}

	// принудительно ограничим время ожидания
	db.SetConnMaxLifetime(5 * time.Second)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0)

	fmt.Println("Calling db.Ping()...")

	err = db.Ping()
	if err != nil {
		fmt.Println("Ping failed with error:", err)
	} else {
		fmt.Println("Ping succeeded (unexpected)")
	}
}
