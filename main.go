package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Shayan-Ghani/GOAD/config"
	"github.com/Shayan-Ghani/GOAD/internal/delivery/cli"
	sqlrepo "github.com/Shayan-Ghani/GOAD/internal/repository/sql"

	_ "github.com/go-sql-driver/mysql"
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

	db := sqlrepo.NewSQLRepo((driver))

	cli := cli.NewCLI(db)

	if err := cli.Exec(os.Args[1:]); err != nil {
		log.Fatalf("%v", err)
	}

}
