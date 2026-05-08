package cache_roble_infrastructure

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MariaDBClient struct {
	db *sql.DB
}

func NewMariaDBClient() (*MariaDBClient, error) {
	dsn := strings.TrimSpace(os.Getenv("CACHE_DB_DSN"))
	if dsn == "" {
		return nil, fmt.Errorf("CACHE_DB_DSN is required")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mariadb: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping mariadb: %w", err)
	}

	return &MariaDBClient{db: db}, nil
}

func (c *MariaDBClient) DB() *sql.DB {
	return c.db
}

func (c *MariaDBClient) Close() error {
	return c.db.Close()
}
