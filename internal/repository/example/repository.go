package example

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domainexample "go-boilerplate/internal/domain/example"
	modelexample "go-boilerplate/internal/model/example"
)

var postgresDialect = goqu.Dialect("postgres")

type ExamplePG struct {
	pool *pgxpool.Pool
}

func NewExamplePG(pool *pgxpool.Pool) *ExamplePG {
	return &ExamplePG{pool: pool}
}

func (r *ExamplePG) List(ctx context.Context, input modelexample.ListInput) ([]*domainexample.Example, error) {
	ds := postgresDialect.
		From(goqu.T("example")).
		Select(goqu.I("id"), goqu.I("name")).
		Order(goqu.I("id").Asc()).
		Limit(uint(input.Limit)).
		Offset(uint(input.Offset))

	query, args, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build list examples query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query examples: %w", err)
	}
	defer rows.Close()

	collectedRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[row])
	if err != nil {
		return nil, fmt.Errorf("collect examples: %w", err)
	}

	items := make([]*domainexample.Example, 0, len(collectedRows))
	for _, dbRow := range collectedRows {
		entity, createErr := domainexample.New(dbRow.ID, dbRow.Name)
		if createErr != nil {
			return nil, fmt.Errorf("create example entity: %w", createErr)
		}

		items = append(items, entity)
	}

	return items, nil
}

func (r *ExamplePG) Count(ctx context.Context) (int, error) {
	ds := postgresDialect.
		From(goqu.T("example")).
		Select(goqu.COUNT(goqu.Star()))

	query, args, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return 0, fmt.Errorf("build count examples query: %w", err)
	}

	var count int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count examples: %w", err)
	}

	return count, nil
}
