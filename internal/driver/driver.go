package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return db, nil
}
