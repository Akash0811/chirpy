package auth

import "testing"

func TestCorrectPassword(t *testing.T) {
	hash, err := HashPassword("VerySecret123")

	match, err := CheckPasswordHash("VerySecret123", hash)
	if err != nil {
		t.Fatalf("Expected no errors but got %v", err)
	}
	if !match {
		t.Fatalf("Expected Password to match")
	}
}

func TestIncorrectPassword(t *testing.T) {
	hash, err := HashPassword("VerySecret123")

	match, err := CheckPasswordHash("AnotherVerySecret456", hash)
	if err != nil {
		t.Fatalf("Expected no errors but got %v", err)
	}
	if match {
		t.Fatalf("Expected Password to not match")
	}
}
