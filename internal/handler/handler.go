package handler

import (
	"github.com/gofiber/fiber/v2"
	"your-project/internal/service"
	
)

type Handler struct {
	cb      *circuitbreaker.CircuitBreaker
	service *service.Service
}

func NewHandler(cb *circuitbreaker.CircuitBreaker, service *service.Service) *Handler {
	return &Handler{
		cb:      cb,
		service: service,
	}
}

func (h *Handler) HandleRequest(c *fiber.Ctx) error {
	result, err := h.cb.Execute(func() (interface{}, error) {
		return h.service.SimulateExternalCall()
	})

	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": err.Error(),
			"state": h.cb.State(),
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
		"state":  h.cb.State(),
	})
}

func (h *Handler) HandleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":       "healthy",
		"circuitState": h.cb.State(),
	})
}
