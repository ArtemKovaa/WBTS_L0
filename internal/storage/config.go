package storage

import (
	"log"
	"os"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(ctx context.Context) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error configuration pool: %v", err)
	}
	config.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Error creating pool: %v", err)
	}
	log.Println("PgPool successfully set up")

	return pool
}
