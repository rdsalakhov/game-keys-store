package server

import (
	"database/sql"
	"github.com/rdsalakhov/game-keys-store/internal/store/mysqlStore"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDb(config.DbConnection)
	if err != nil {
		return err
	}
	defer db.Close()
	store := mysqlStore.New(db)
	server := NewServer(store)
	return http.ListenAndServe(config.Port, server)
}

func newDb(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
