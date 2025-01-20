package main

import (
	"./config"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sony/gobreaker"
	"time"
)

// Circuit Breaker konfigürasyonu
var cb *gobreaker.CircuitBreaker

func init() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Config yüklenemedi: %v", err))
	}

	settings := gobreaker.Settings{
		Name:        cfg.CircuitBreaker.Name,
		MaxRequests: cfg.CircuitBreaker.MaxRequests,
		Interval:    time.Duration(cfg.CircuitBreaker.Interval) * time.Second,
		Timeout:     time.Duration(cfg.CircuitBreaker.Timeout) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= cfg.CircuitBreaker.MinimumRequests && failureRatio >= cfg.CircuitBreaker.FailureRatio
		},
	}

	cb = gobreaker.NewCircuitBreaker(settings)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Config yüklenemedi: %v", err))
	}

	app := fiber.New()

	app.Get("/api", handleRequest)
	app.Get("/health", handleHealth)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	app.Listen(addr)
}

func handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":       "healthy",
		"circuitState": cb.State().String(),
	})
}

func handleRequest(c *fiber.Ctx) error {
	result, err := cb.Execute(func() (interface{}, error) {
		// Simüle edilmiş dış servis çağrısı
		if time.Now().Unix()%2 == 0 { // Rastgele hata üretimi
			return nil, errors.New("service error")
		}
		return "Success!", nil
	})

	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": err.Error(),
			"state": cb.State().String(),
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
		"state":  cb.State().String(),
	})
}
