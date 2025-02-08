package main

import (
	"database/sql"
	"log"

	"github.com/Shayan-Ghani/GOAD/config"
	"github.com/Shayan-Ghani/GOAD/internal/delivery/tcpserver"
	sqlrepository "github.com/Shayan-Ghani/GOAD/internal/repository/sql"

	_ "github.com/go-sql-driver/mysql"

)

func main() {
	driver, err := sql.Open("mysql", config.DefaultDSN)

	if err != nil {
		log.Fatalln(err)
	}

	ping := driver.Ping()
	if ping != nil {
		log.Fatal(ping)
	}

	repo := sqlrepository.NewSQLRepo((driver))

    server, err := tcpserver.NewServer(repo)

    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    
    if err := server.Start(); err != nil {
        log.Fatalf("Server error: %v", err)
    }

}