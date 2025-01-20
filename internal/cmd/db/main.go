package cmd

import (
	"database/sql"
	"fmt"
	"gocasts/ToDoApp/internal/db"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	dsn = "root:root@tcp(localhost:3306)/todo_dev?multiStatements=true"
)

func Migrate() {
	driver, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	ping := driver.Ping()
	if ping != nil {
		log.Fatal(ping)
	}
	
	_, b, _, _ := runtime.Caller(0)
	
	
	currentDir := filepath.Dir(b)
	migrationPath := filepath.Join(currentDir, "internal", "db", "migrations")
	
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