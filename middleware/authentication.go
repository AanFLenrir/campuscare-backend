package middleware

import (
	"github.com/AanFLenrir/PetrivaBackend/config"
	"github.com/AanFLenrir/PetrivaBackend/helper"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// AuthMiddleware validates JWT and stores user info in context locals.
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.Error(c, fiber.StatusUnauthorized, "Missing Authorization header")
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return helper.Error(c, fiber.StatusUnauthorized, "Invalid Authorization header format")
		}
		tokenStr := parts[1]
		claims, err := helper.ValidateToken(cfg, tokenStr)
		if err != nil {
			return helper.Error(c, fiber.StatusUnauthorized, "Invalid token")
		}
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}
