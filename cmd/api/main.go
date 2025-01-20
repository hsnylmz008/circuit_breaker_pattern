package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"your-project/internal/config"
	"your-project/internal/handler"
	"your-project/internal/router"
	"your-project/internal/service"
	"your-project/pkg/circuitbreaker"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	// Initialize circuit breaker
	cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.Config{
		Name:            cfg.CircuitBreaker.Name,
		MaxRequests:     cfg.CircuitBreaker.MaxRequests,
		Interval:        time.Duration(cfg.CircuitBreaker.Interval) * time.Second,
		Timeout:         time.Duration(cfg.CircuitBreaker.Timeout) * time.Second,
		FailureRatio:    cfg.CircuitBreaker.FailureRatio,
		MinimumRequests: cfg.CircuitBreaker.MinimumRequests,
	})

	// Initialize service
	svc := service.NewService()

	// Initialize handler
	h := handler.NewHandler(cb, svc)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	router.SetupRoutes(app, h)

	// Graceful shutdown için
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
	}()

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
