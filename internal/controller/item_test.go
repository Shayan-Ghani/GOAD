package controller_test

import (
	"database/sql"
	"fmt"
	"gocasts/ToDoApp/internal/controller"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dsn = "root:root@tcp(localhost:3306)/todo_test"
)
var(
	driver, _ = sql.Open("mysql", dsn)
	c = controller.NewSQLTodoController(driver)
)

type result struct {
	got string
	want string
}


func TestAddItem(t *testing.T) {
	var r = &result{}
	ts , err := time.Parse("2006-01-02 15:04:05", "2025-03-05 15:05:10");

	if err != nil {
		t.Fatalf("got %q, wanted %q", err, "2025-03-05 15:05:10 +0000 UTC")
	}
	fmt.Printf("time.Now().String(): %v\n", ts)

	var ti = time.Time{}
	if err := c.AddItem("test", "test", ti); err != nil {
		r.got = err.Error()
	} 
	if r.got != r.want{
		t.Errorf("got %q, wanted %q", r.got, r.want)
	}

}
