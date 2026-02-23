package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExamplePG struct {
	pool *pgxpool.Pool
}

func NewExamplePG(pool *pgxpool.Pool) *ExamplePG {
	return &ExamplePG{pool: pool}
}

func (r *ExamplePG) CountExamples(ctx context.Context) (int64, error) {
	var count int64

	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM example`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count examples: %w", err)
	}

	return count, nil
}
