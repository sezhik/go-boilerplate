package example_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	httpexample "go-boilerplate/internal/handlers/http/example"
	modelexample "go-boilerplate/internal/model/example"
	presenterexample "go-boilerplate/internal/presenter/example"
)

func TestHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		target       string
		prepare      func(usecase *Mockusecase, presenter *Mockpresenter)
		expectations func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "usecase failed",
			target: "/api/examples",
			prepare: func(usecase *Mockusecase, presenter *Mockpresenter) {
				usecase.EXPECT().
					List(gomock.Any(), modelexample.ListInput{Limit: 10, Offset: 0}).
					Return(modelexample.ListOutput{}, assert.AnError)
			},
			expectations: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
				assert.JSONEq(t, `{"error":"failed to list examples"}`, recorder.Body.String())
			},
		},
		{
			name:   "returns presented page",
			target: "/api/examples?limit=2&offset=4",
			prepare: func(usecase *Mockusecase, presenter *Mockpresenter) {
				output := modelexample.ListOutput{
					Items:  []modelexample.Item{{ID: 1, Name: "first"}},
					Limit:  2,
					Offset: 4,
					Total:  5,
				}

				usecase.EXPECT().
					List(gomock.Any(), modelexample.ListInput{Limit: 2, Offset: 4}).
					Return(output, nil)
				presenter.EXPECT().
					List(output).
					Return(presenterexample.ListResponse{
						Items:  []presenterexample.ItemResponse{{ID: 1, Name: "first"}},
						Limit:  2,
						Offset: 4,
						Total:  5,
					})
			},
			expectations: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				assert.Equal(t, http.StatusOK, recorder.Code)
				assert.JSONEq(t, `{"items":[{"id":1,"name":"first"}],"limit":2,"offset":4,"total":5}`, recorder.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usecase := NewMockusecase(ctrl)
			presenter := NewMockpresenter(ctrl)
			tt.prepare(usecase, presenter)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodGet, tt.target, nil)

			handler := httpexample.NewHandler(usecase, presenter, slog.New(slog.DiscardHandler))
			handler.List(ctx)

			tt.expectations(t, recorder)
		})
	}
}
