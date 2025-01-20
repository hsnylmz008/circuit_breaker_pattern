package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	cacheStore = cache.New(5*time.Minute, 10*time.Minute)
)

func Cache(duration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// GET istekleri için cache
		if c.Method() != "GET" {
			return c.Next()
		}

		key := c.Path()
		if val, found := cacheStore.Get(key); found {
			return c.JSON(val)
		}

		err := c.Next()
		if err != nil {
			return err
		}

		cacheStore.Set(key, c.Response(), duration)
		return nil
	}
}
