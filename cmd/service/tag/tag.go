package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Shayan-Ghani/GOAD/config"
	sqlrepository "github.com/Shayan-Ghani/GOAD/internal/repository/sql"
	tagsvc "github.com/Shayan-Ghani/GOAD/internal/service/tag"
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

	repo := sqlrepository.NewTagRepo((driver))

	r := http.NewServeMux()
	s := tagsvc.NewTagService(repo)
	tagsvc.Handle(r, s)

	fmt.Println("listening on 8788 ...")
	http.ListenAndServe(":8788", r)
}
