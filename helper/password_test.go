package helper_test

import (
	"testing"

	"github.com/AanFLenrir/PetrivaBackend/helper"
)

func TestHashPassword_Success(t *testing.T) {
	password := "rahasia123"
	hash, err := helper.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if hash == "" || hash == password {
		t.Fatalf("expected valid bcrypt hash, got: %s", hash)
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	password := "rahasia123"
	hash, err := helper.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	err = helper.CheckPassword(hash, password)
	if err != nil {
		t.Fatalf("expected password match, got error: %v", err)
	}
}

func TestCheckPassword_Incorrect(t *testing.T) {
	password := "rahasia123"
	hash, err := helper.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	err = helper.CheckPassword(hash, "passwordsalah")
	if err == nil {
		t.Fatalf("expected error for mismatched password, got nil")
	}
}
