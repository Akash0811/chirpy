package auth

import (
	"crypto/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCorrectToken(t *testing.T) {
	id := uuid.New()
	tokenSecret := rand.Text()

	duration, err := time.ParseDuration("120s")
	if err != nil {
		t.Fatalf("Expected string to be parse as time but got %v", err)
	}

	tokenString, err := MakeJWT(id, tokenSecret, duration)
	if err != nil {
		t.Fatalf("Expected MakeJWT to work but got %v", err)
	}

	gotId, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("Expected ValidateJWT to work but got %v", err)
	}

	if gotId != id {
		t.Fatal("Expected ids to match")
	}
}

func TestDifferentSecret(t *testing.T) {
	id := uuid.New()
	tokenSecret1 := rand.Text()
	tokenSecret2 := rand.Text()

	if tokenSecret1 == tokenSecret2 {
		t.Fatalf("Expected secrets to be different")
	}

	duration, err := time.ParseDuration("120s")
	if err != nil {
		t.Fatalf("Expected string to be parse as time but got %v", err)
	}

	tokenString, err := MakeJWT(id, tokenSecret1, duration)
	if err != nil {
		t.Fatalf("Expected MakeJWT to work but got %v", err)
	}

	_, err = ValidateJWT(tokenString, tokenSecret2)
	if err == nil {
		t.Fatalf("Expected ValidateJWT to fail")
	}
}

func TestSecretTimeElapsed(t *testing.T) {
	id := uuid.New()
	tokenSecret := rand.Text()

	duration, err := time.ParseDuration("5s")
	if err != nil {
		t.Fatalf("Expected string to be parse as time but got %v", err)
	}
	moreDuration, err := time.ParseDuration("7s")
	if err != nil {
		t.Fatalf("Expected string to be parse as time but got %v", err)
	}

	tokenString, err := MakeJWT(id, tokenSecret, duration)
	if err != nil {
		t.Fatalf("Expected MakeJWT to work but got %v", err)
	}

	time.Sleep(moreDuration)

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatalf("Expected ValidateJWT to fail")
	}
}

func TestBearerPresent(t *testing.T) {
	headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer  Passme  "},
	}
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("Expected no errors but got %v", err)
	}
	if token != "Passme" {
		t.Fatalf("Expected token: Passme but got %s", token)
	}
}

func TestBearerNotPresent(t *testing.T) {
	headers := http.Header{
		"Content-Type":   []string{"application/json"},
		"Authentication": []string{"Bear  Passme  "},
	}
	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("Expected errors but got token")
	}
}
