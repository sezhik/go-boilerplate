package example

import (
	"context"

	modelexample "go-boilerplate/internal/model/example"
	presenterexample "go-boilerplate/internal/presenter/example"
)

//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go

type usecase interface {
	List(ctx context.Context, input modelexample.ListInput) (modelexample.ListOutput, error)
}

type presenter interface {
	List(output modelexample.ListOutput) presenterexample.ListResponse
}
