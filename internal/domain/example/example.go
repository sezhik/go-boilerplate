package example

import (
	"fmt"
	"strings"

	"go-boilerplate/internal/domain/domainerror"
)

type Example struct {
	id   int64
	name string
}

func New(id int64, name string) (*Example, error) {
	normalizedName := strings.TrimSpace(name)

	if id <= 0 {
		return nil, fmt.Errorf("example id must be positive: %w", domainerror.Invalid)
	}
	if normalizedName == "" {
		return nil, fmt.Errorf("example name is required: %w", domainerror.Invalid)
	}

	return &Example{id: id, name: normalizedName}, nil
}

func (e *Example) ID() int64 {
	return e.id
}

func (e *Example) Name() string {
	return e.name
}
