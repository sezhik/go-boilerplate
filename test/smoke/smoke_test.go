//go:build smoke

// Package smoke runs HTTP smoke checks against a live dev stack. Opt-in via
// `-tags=smoke` so that `go test ./...` never hits the network.
package smoke

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

const defaultBaseURL = "http://localhost:8080"

func baseURL(t *testing.T) string {
	t.Helper()

	// .env in repo root is the source of truth for local dev; existing env
	// vars take priority.
	_ = godotenv.Load("../../.env")

	if url := os.Getenv("SMOKE_BASE_URL"); url != "" {
		return url
	}

	return defaultBaseURL
}

func get(t *testing.T, url string) (int, []byte) {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, body
}

func TestSmoke(t *testing.T) {
	base := baseURL(t)

	t.Run("health returns ok", func(t *testing.T) {
		code, body := get(t, base+"/_health")
		require.Equal(t, http.StatusOK, code)
		require.JSONEq(t, `{"status":"ok","service":"go-boilerplate"}`, string(body))
	})

	t.Run("examples list returns page", func(t *testing.T) {
		code, body := get(t, base+"/api/examples")
		require.Equalf(t, http.StatusOK, code, "list examples failed: %s", string(body))

		var result struct {
			Items  []map[string]any `json:"items"`
			Limit  int              `json:"limit"`
			Offset int              `json:"offset"`
			Total  int              `json:"total"`
		}
		require.NoErrorf(t, json.Unmarshal(body, &result), "decode body: %s", string(body))
		require.Equal(t, 10, result.Limit)
		require.Equal(t, 0, result.Offset)
	})
}
