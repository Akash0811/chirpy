package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	bearerString := headers.Get("Authorization")
	if bearerString == "" {
		return "", fmt.Errorf("Authorization Header not present")
	}

	return strings.TrimSpace(strings.TrimPrefix(bearerString, "ApiKey")), nil
}
