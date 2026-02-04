package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword returned empty string")
	}

	if hash == password {
		t.Fatal("HashPassword should not return the plain password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	t.Run("correct password", func(t *testing.T) {
		match, err := CheckPasswordHash(password, hash)
		if err != nil {
			t.Fatalf("CheckPasswordHash failed: %v", err)
		}
		if !match {
			t.Error("CheckPasswordHash should return true for correct password")
		}
	})

	t.Run("incorrect password", func(t *testing.T) {
		wrongPassword := "wrongpassword"
		match, err := CheckPasswordHash(wrongPassword, hash)
		if err != nil {
			t.Fatalf("CheckPasswordHash failed: %v", err)
		}
		if match {
			t.Error("CheckPasswordHash should return false for incorrect password")
		}
	})
}

func TestHashPassword_UniqueHashes(t *testing.T) {
	password := "samepassword"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword should produce unique hashes for the same password (due to salt)")
	}
}
