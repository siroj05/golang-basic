package golang_database

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestEmpty(t *testing.T) {

}

func TestOpenConnection(t *testing.T) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:8080)/golang_database")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	// gunakan db jangan lupa di close

}
