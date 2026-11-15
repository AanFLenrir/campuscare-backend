package helper_test

import (
	"testing"
	"time"

	"github.com/AanFLenrir/campuscare-backend/config"
	"github.com/AanFLenrir/campuscare-backend/helper"
	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken_Success(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "secret-key-campuscare-test-12345",
		JWTExpireHour: 1,
	}

	userID := "usr-uuid-001"
	role := "MAHASISWA"

	tokenStr, err := helper.GenerateToken(cfg, userID, role)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}
	if tokenStr == "" {
		t.Fatalf("expected non-empty token string")
	}

	claims, err := helper.ValidateToken(cfg, tokenStr)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
}

func TestValidateToken_InvalidSecret(t *testing.T) {
	cfg1 := &config.Config{
		JWTSecret:     "secret-key-1",
		JWTExpireHour: 1,
	}
	cfg2 := &config.Config{
		JWTSecret:     "secret-key-2",
		JWTExpireHour: 1,
	}

	tokenStr, err := helper.GenerateToken(cfg1, "usr-001", "ADMIN")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = helper.ValidateToken(cfg2, tokenStr)
	if err == nil {
		t.Fatalf("expected error when validating token with wrong secret, got nil")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "secret-key-expired-test",
		JWTExpireHour: 1,
	}

	// Buat token yang sudah expired secara manual
	claims := helper.Claims{
		UserID: "usr-001",
		Role:   "PETUGAS",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = helper.ValidateToken(cfg, tokenStr)
	if err == nil {
		t.Fatalf("expected error validating expired token, got nil")
	}
}
