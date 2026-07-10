//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package example

import (
	"context"

	domainexample "go-boilerplate/internal/domain/example"
	modelexample "go-boilerplate/internal/model/example"
)

type examples interface {
	List(ctx context.Context, input modelexample.ListInput) ([]*domainexample.Example, error)
	Count(ctx context.Context) (int, error)
}
