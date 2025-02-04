package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Shayan-Ghani/GOAD/config"
	"github.com/Shayan-Ghani/GOAD/internal/db"
)

func main() {
	driver, err := sql.Open("mysql", config.DefaultDSN)
	if err != nil {
		log.Fatal(err)
	}

	ping := driver.Ping()
	if ping != nil {
		log.Fatal(ping)
	}

	_, b, _, _ := runtime.Caller(0)

	currentDir := filepath.Dir(b)
	migrationPath := filepath.Join(currentDir, "..", "internal", "db", "migrations")
	fmt.Println(migrationPath)

	dbManager, err := db.NewDBManager(driver, migrationPath)
	if err != nil {
		log.Fatal(err)
	}

	switch action := os.Getenv("MIGRATE"); action {
	case "UP":
		if err := dbManager.MigrateUp(); err != nil {
			log.Fatal(err)
		}
	case "DOWN":
		if err := dbManager.MigrateDown(); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("** No Migration **")
	}
}
