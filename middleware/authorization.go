package middleware

import (
	"github.com/AanFLenrir/PetrivaBackend/helper"
	"github.com/gofiber/fiber/v2"
)

// Authorize checks if the role stored in context locals is allowed.
func Authorize(allowedRoles ...string) fiber.Handler {
	roleSet := make(map[string]struct{})
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		val := c.Locals("role")
		r, ok := val.(string)
		if !ok {
			return helper.Error(c, fiber.StatusForbidden, "Role not found in context")
		}
		if _, ok := roleSet[r]; !ok {
			return helper.Error(c, fiber.StatusForbidden, "Access denied for role")
		}
		return c.Next()
	}
}
