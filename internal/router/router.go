package router

import (
	"github.com/gofiber/adaptor"

"github.com/gofiber/fiber/v2"
	"github
	"github.com/gofiber/fiber/v2"
	"time"
	"your-project/internal/handler"
	"golang.org/x/time/rate"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, h *handler.Handler) {
	// Middleware'ler

	
	rateLimiter := middleware.NewRateLimiter(rate.Limit(100), 10)
	app.Use(rateLimiter.Middleware())

	// API endpoints
	api := app.Group("/api")

	
	// Cache'li endpoint örneği

	
	// Health check

	
	// Metrics endpoint
}

} 