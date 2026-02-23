package configure

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-boilerplate/internal/util/env"
)

func MustNewDB() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pool, err := pgxpool.New(ctx, fmt.Sprintf("host=pgbouncer dbname=%s port=5432 user=%s password=%s", env.MustGet("PG_DB"), env.MustGet("PG_USER"), env.MustGet("PG_PASS")))
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var (
		okPing       bool
		okMigrations bool
	)

	for {
		if okMigrations {
			break
		}

		select {
		case <-ctx.Done():
			if !okPing {
				panic("failed to ping database")
			}
			panic("timeout waiting for the db migrations")
		case <-ticker.C:
			if !okPing {
				err = pool.Ping(ctx)
				okPing = err == nil
			}

			if okPing {
				err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'schema_migrations')").Scan(&okMigrations)
				if err != nil {
					panic("failed to check migrations status")
				}
			}
		}
	}

	return pool
}
