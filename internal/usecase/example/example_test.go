package example_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	domainexample "go-boilerplate/internal/domain/example"
	modelexample "go-boilerplate/internal/model/example"
	usecaseexample "go-boilerplate/internal/usecase/example"
)

func TestUsecaseList(t *testing.T) {
	t.Parallel()

	newEntity := func(t *testing.T, id int64, name string) *domainexample.Example {
		t.Helper()

		entity, err := domainexample.New(id, name)
		require.NoError(t, err)

		return entity
	}

	tests := []struct {
		name         string
		input        modelexample.ListInput
		prepare      func(t *testing.T, examples *Mockexamples)
		expectations func(t *testing.T, output modelexample.ListOutput, err error)
	}{
		{
			name:  "repository failed",
			input: modelexample.ListInput{},
			prepare: func(t *testing.T, examples *Mockexamples) {
				t.Helper()

				examples.EXPECT().
					List(gomock.Any(), modelexample.ListInput{Limit: 10, Offset: 0}).
					Return(nil, assert.AnError)
			},
			expectations: func(t *testing.T, output modelexample.ListOutput, err error) {
				t.Helper()

				assert.ErrorIs(t, err, assert.AnError)
				assert.Equal(t, modelexample.ListOutput{}, output)
			},
		},
		{
			name:  "count failed",
			input: modelexample.ListInput{},
			prepare: func(t *testing.T, examples *Mockexamples) {
				t.Helper()

				examples.EXPECT().
					List(gomock.Any(), modelexample.ListInput{Limit: 10, Offset: 0}).
					Return([]*domainexample.Example{newEntity(t, 1, "first")}, nil)
				examples.EXPECT().
					Count(gomock.Any()).
					Return(0, assert.AnError)
			},
			expectations: func(t *testing.T, output modelexample.ListOutput, err error) {
				t.Helper()

				assert.ErrorIs(t, err, assert.AnError)
				assert.Equal(t, modelexample.ListOutput{}, output)
			},
		},
		{
			name:  "defaults pagination and returns page",
			input: modelexample.ListInput{Offset: -7},
			prepare: func(t *testing.T, examples *Mockexamples) {
				t.Helper()

				examples.EXPECT().
					List(gomock.Any(), modelexample.ListInput{Limit: 10, Offset: 0}).
					Return([]*domainexample.Example{newEntity(t, 1, "first")}, nil)
				examples.EXPECT().
					Count(gomock.Any()).
					Return(1, nil)
			},
			expectations: func(t *testing.T, output modelexample.ListOutput, err error) {
				t.Helper()

				assert.NoError(t, err)
				assert.Equal(t, modelexample.ListOutput{
					Items:  []modelexample.Item{{ID: 1, Name: "first"}},
					Limit:  10,
					Offset: 0,
					Total:  1,
				}, output)
			},
		},
		{
			name:  "keeps custom pagination",
			input: modelexample.ListInput{Limit: 25, Offset: 10},
			prepare: func(t *testing.T, examples *Mockexamples) {
				t.Helper()

				examples.EXPECT().
					List(gomock.Any(), modelexample.ListInput{Limit: 25, Offset: 10}).
					Return(nil, nil)
				examples.EXPECT().
					Count(gomock.Any()).
					Return(42, nil)
			},
			expectations: func(t *testing.T, output modelexample.ListOutput, err error) {
				t.Helper()

				assert.NoError(t, err)
				assert.Equal(t, modelexample.ListOutput{
					Items:  []modelexample.Item{},
					Limit:  25,
					Offset: 10,
					Total:  42,
				}, output)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			examples := NewMockexamples(ctrl)
			tt.prepare(t, examples)

			u := usecaseexample.New(examples)
			output, err := u.List(context.Background(), tt.input)

			tt.expectations(t, output, err)
		})
	}
}
