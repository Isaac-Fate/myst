package cmd_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateDatabase(t *testing.T) {
	db, err := sql.Open("sqlite3", "./test.db")

	if err != nil {
		t.Error(err)
	}

	// SQL statement to create the todos table if it doesn't exist
	sqlStmt := `
	 CREATE TABLE IF NOT EXISTS todos (
	  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	  title TEXT
	 );`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		t.Error(err)
	}

	db.Close()
}
