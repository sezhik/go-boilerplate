package example

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-boilerplate/internal/handlers/http/httputil"
	modelexample "go-boilerplate/internal/model/example"
)

type Handler struct {
	usecase   usecase
	presenter presenter
	logger    *slog.Logger
}

func NewHandler(usecase usecase, presenter presenter, logger *slog.Logger) *Handler {
	return &Handler{usecase: usecase, presenter: presenter, logger: logger}
}

func (h *Handler) List(c *gin.Context) {
	limit := httputil.ParseIntOrDefault(c.Query("limit"), 10)
	offset := httputil.ParseIntOrDefault(c.Query("offset"), 0)

	result, err := h.usecase.List(c.Request.Context(), modelexample.ListInput{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "list examples failed", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list examples"})
		return
	}

	c.JSON(http.StatusOK, h.presenter.List(result))
}
