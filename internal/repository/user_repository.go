package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	db "user-api/db/sqlc"
)

func NewRepository(conn *pgx.Conn) db.Querier {
	return db.New(conn)
}

func RunMigration(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id   SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			dob  DATE NOT NULL
		);
	`)
	return err
}
