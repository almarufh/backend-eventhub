package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	User     string
	Password string
	Host     string
	Port     string
	Db_name  string
}

func NewPostgreSQL() *Postgres {
	return &Postgres{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Db_name:  os.Getenv("POSTGRES_DB"),
	}
}

func (p *Postgres) Connect() (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.User, p.Password, p.Host, p.Port, p.Db_name)
	return pgxpool.New(context.Background(), url)
}
