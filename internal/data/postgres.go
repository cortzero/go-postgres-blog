package data

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

const SQL_SCHEMA_URL = "./database/schema.sql"

func getConnection() (*sql.DB, error) {
	uri := os.Getenv("DATABASE_URI")
	return sql.Open("postgres", uri)
}

func MakeMigration(db *sql.DB) error {
	data, err := os.ReadFile(SQL_SCHEMA_URL)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(data))
	if err != nil {
		return err
	} else {
		return nil
	}
}
