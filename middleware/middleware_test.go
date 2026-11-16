package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AanFLenrir/PetrivaBackend/config"
	"github.com/AanFLenrir/PetrivaBackend/helper"
	"github.com/AanFLenrir/PetrivaBackend/middleware"
	"github.com/gofiber/fiber/v2"
)

func setupTestApp(cfg *config.Config) *fiber.App {
	app := fiber.New()

	// Protected route with AuthMiddleware
	app.Get("/protected", middleware.AuthMiddleware(cfg), func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		role := c.Locals("role").(string)
		return c.JSON(fiber.Map{
			"message": "success",
			"userID":  userID,
			"role":    role,
		})
	})

	// Role-restricted route with AuthMiddleware + Authorize("ADMIN")
	app.Get("/admin-only", middleware.AuthMiddleware(cfg), middleware.Authorize("ADMIN"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "welcome admin"})
	})

	// Multiple allowed roles route
	app.Get("/staff-only", middleware.AuthMiddleware(cfg), middleware.Authorize("ADMIN", "VETERINARIAN"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "welcome clinic staff"})
	})

	return app
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpireHour: 1}
	app := setupTestApp(cfg)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpireHour: 1}
	app := setupTestApp(cfg)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic 123456")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpireHour: 1}
	app := setupTestApp(cfg)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthMiddleware_Success(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpireHour: 1}
	app := setupTestApp(cfg)

	token, err := helper.GenerateToken(cfg, "usr-123", "PET_OWNER")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status %d, got %d. body: %s", fiber.StatusOK, resp.StatusCode, string(body))
	}
}

func TestAuthorize_AdminAllowed(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpireHour: 1}
	app := setupTestApp(cfg)

	token, err := helper.GenerateToken(cfg, "admin-1", "ADMIN")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func TestAuthorize_ForbiddenRole(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpireHour: 1}
	app := setupTestApp(cfg)

	// PET_OWNER trying to access /admin-only
	token, err := helper.GenerateToken(cfg, "owner-1", "PET_OWNER")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Errorf("expected status %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}
}

func TestAuthorize_MultipleAllowedRoles(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret", JWTExpireHour: 1}
	app := setupTestApp(cfg)

	// Veterinarian accessing /staff-only (allowed: ADMIN, VETERINARIAN)
	tokenVet, _ := helper.GenerateToken(cfg, "vet-1", "VETERINARIAN")
	req := httptest.NewRequest(http.MethodGet, "/staff-only", nil)
	req.Header.Set("Authorization", "Bearer "+tokenVet)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	// PET_OWNER accessing /staff-only (forbidden)
	tokenOwner, _ := helper.GenerateToken(cfg, "owner-1", "PET_OWNER")
	req2 := httptest.NewRequest(http.MethodGet, "/staff-only", nil)
	req2.Header.Set("Authorization", "Bearer "+tokenOwner)
	resp2, err := app.Test(req2, -1)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != fiber.StatusForbidden {
		t.Errorf("expected status %d, got %d", fiber.StatusForbidden, resp2.StatusCode)
	}
}
