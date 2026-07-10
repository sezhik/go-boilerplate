package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"go-boilerplate/internal/handlers/http/health"
)

func TestHandlerHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name  string
		input struct {
			method string
			target string
		}
		expectations func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "returns healthy service response",
			input: struct {
				method string
				target string
			}{
				method: http.MethodGet,
				target: "/_health",
			},
			expectations: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()

				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
				require.JSONEq(t, `{"status":"ok","service":"go-boilerplate"}`, recorder.Body.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(tt.input.method, tt.input.target, nil)

			handler := health.NewHandler()
			handler.Handle(ctx)

			tt.expectations(t, recorder)
		})
	}
}
