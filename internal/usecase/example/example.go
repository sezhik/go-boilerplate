package example

import (
	"context"
	"fmt"

	modelexample "go-boilerplate/internal/model/example"
)

type Usecase struct {
	examples examples
}

func New(examples examples) *Usecase {
	return &Usecase{examples: examples}
}

func (u *Usecase) List(ctx context.Context, input modelexample.ListInput) (modelexample.ListOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := max(input.Offset, 0)

	entities, err := u.examples.List(ctx, modelexample.ListInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return modelexample.ListOutput{}, fmt.Errorf("list examples: %w", err)
	}

	total, err := u.examples.Count(ctx)
	if err != nil {
		return modelexample.ListOutput{}, fmt.Errorf("count examples: %w", err)
	}

	items := make([]modelexample.Item, 0, len(entities))
	for _, entity := range entities {
		items = append(items, modelexample.Item{
			ID:   entity.ID(),
			Name: entity.Name(),
		})
	}

	return modelexample.ListOutput{
		Items:  items,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}
