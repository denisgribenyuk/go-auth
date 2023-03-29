package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

var DB *sqlx.DB

func InitDB(cfg Config, maxOpenConnections int, maxIdleConnections int) {
	DB = NewPostgresDB(cfg, maxOpenConnections, maxIdleConnections)
}

// NewPostgresDB Creates and returns pointer to new DB instance.
func NewPostgresDB(cfg Config, maxOpenConnections int, maxIdleConnections int) *sqlx.DB {
	db := sqlx.MustOpen("pgx", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password))

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
