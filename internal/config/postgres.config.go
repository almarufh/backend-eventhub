package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	User     string
	Password string
	Host     string
	Port     string
	Db_name  string
}

func NewPostgreSQL(user, password, host, port, dbName string) *Postgres {
	return &Postgres{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Db_name:  dbName,
	}
}

func (p *Postgres) Connect() *pgxpool.Pool {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.User, p.Password, p.Host, p.Port, p.Db_name)
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Printf("Gagal inisialisasi database pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Printf("Database tidak merespon: %v", err)
	}

	log.Printf("Connected PostgreSQL successfully")
	return pool
}
