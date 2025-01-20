package main

import (
	"database/sql"
	"gocasts/ToDoApp/internal/controller"
	"gocasts/ToDoApp/internal/delivery/cli"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dsn = "root:root@tcp(localhost:3306)/todo_dev"
)

func main() {
	driver, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	ping := driver.Ping()
	if ping != nil {
		log.Fatal(ping)
	}

	db := controller.NewSQLTodoController(driver)

	cli := cli.NewCLI(db)


	if err := cli.Exec(os.Args[1:]); err != nil {
		log.Fatalf("%v", err)
	}
	
}
