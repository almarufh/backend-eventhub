package pkg

import (
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound       = errors.New("data not found")
	ErrConflict       = errors.New("data already exists")
	ErrEmailExists    = errors.New("email already exists")
	ErrAuthExists     = errors.New("you have token access active")
	ErrForeignKeyFail = errors.New("referenced record not found")
	ErrCheckViolation = errors.New("invalid input value violates constraint")
	ErrNotNullFail    = errors.New("required field cannot be empty")
)

func ParseError(err error) error {
	if err == nil {
		return nil
	}

	// 1. Cek pgx ErrNoRows (SELECT / QueryRow kosong)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Println(ErrNotFound)
		return ErrNotFound
	}

	// 2. Cek Postgres error spesifik dari pgconn
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			switch pgErr.ConstraintName {
			case "users_email_key":
				log.Println(ErrEmailExists)
				return ErrEmailExists
			case "auth_users_token_key":
				log.Println(ErrAuthExists)
				return ErrAuthExists
			default:
				log.Printf("%s: %s", ErrConflict, pgErr.ConstraintName)
				return fmt.Errorf("%s: %s", ErrConflict, pgErr.ConstraintName)
			}

		case "23503": // foreign_key_violation
			log.Printf("%s: %s", ErrForeignKeyFail, pgErr.ConstraintName)
			return fmt.Errorf("%s: %s", ErrForeignKeyFail, pgErr.ConstraintName)

		case "23514": // check_violation
			log.Printf("%s: %s", ErrCheckViolation, pgErr.ConstraintName)
			return fmt.Errorf("%s: %s", ErrCheckViolation, pgErr.ConstraintName)

		case "23502": // not_null_violation
			log.Printf("%s: column %s", ErrNotNullFail, pgErr.ColumnName)
			return fmt.Errorf("%s: column %s", ErrNotNullFail, pgErr.ColumnName)
		}
	}

	// 3. Fallback jika bukan error constraint database yang dipetakan
	return err
}
