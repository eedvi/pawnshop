package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler exposes Prometheus metrics endpoint
type MetricsHandler struct{}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// Handler returns the Prometheus metrics handler
func (h *MetricsHandler) Handler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}

// RegisterRoutes registers the metrics endpoint
func (h *MetricsHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/metrics", h.Handler())
}
