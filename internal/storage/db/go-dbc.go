package db

import (
	"URLShortener/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func NewDBConn(cfg *config.AppConfig) (*bun.DB, error) {
	dbType, dsn := parseConfig(cfg)

	switch dbType {
	case "postgres":
		return newPostgresDBConn(dsn), nil
	case "mysql":
		return newMySQLDBConn(dsn), nil
	case "sqlite3":
		return newSQLiteDBConn(dsn), nil
	default:
		log.Fatalf("Failed to initialize database. Unknown type: %s", dbType)
		return nil, nil
	}
}

// parseConfig etxracts the configuration info and builds a database connection string to create the connection
func parseConfig(cfg *config.AppConfig) (string, string) {
	switch cfg.DBType {
	case "postgres":
		return "postgres", fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
		)
	case "mysql":
		return "mysql", fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
		)
	case "sqlite":
		return "sqlite3", cfg.DBPath // SQLite uses file-based DSN
	default:
		log.Fatalf("Unsupported database type: %s", cfg.DBType)
		return "", ""
	}
}

// newPostgresDBConn Create a postgres connection
func newPostgresDBConn(dsn string) *bun.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return bun.NewDB(db, nil)
}

// newMySQLDBConn Create a MySQL connection
func newMySQLDBConn(dsn string) *bun.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return bun.NewDB(db, nil)
}

// newSQLiteDBConn Create a SQLite connection
func newSQLiteDBConn(dsn string) *bun.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return bun.NewDB(db, sqlitedialect.New())
}
