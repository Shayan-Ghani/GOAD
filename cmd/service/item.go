package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Shayan-Ghani/GOAD/config"
	sqlrepository "github.com/Shayan-Ghani/GOAD/internal/repository/sql"
	itemsvc "github.com/Shayan-Ghani/GOAD/internal/service/item"
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

	repo := sqlrepository.NewItemRepo((driver))

	r := http.NewServeMux()
	s := itemsvc.NewItemService(repo, config.TagSvcAddr)
	itemsvc.Handle(r, &s)

	fmt.Println("listening on 8787 ...")
	http.ListenAndServe(":8787", r)
}
