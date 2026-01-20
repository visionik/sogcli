package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/zalando/go-keyring"
)

const serviceName = "sog"

// SetPassword stores a password in the system keyring.
func SetPassword(email, password string) error {
	return keyring.Set(serviceName, email, password)
}

// GetPassword retrieves a password from the system keyring.
// Falls back to environment variable SOG_PASSWORD_<email> if keyring fails.
func GetPassword(email string) (string, error) {
	// Try keyring first
	password, err := keyring.Get(serviceName, email)
	if err == nil {
		return password, nil
	}

	// Fall back to environment variable
	envKey := "SOG_PASSWORD_" + sanitizeEnvKey(email)
	if envPass := os.Getenv(envKey); envPass != "" {
		return envPass, nil
	}

	return "", fmt.Errorf("password not found for %s (tried keyring and %s)", email, envKey)
}

// DeletePassword removes a password from the system keyring.
func DeletePassword(email string) error {
	return keyring.Delete(serviceName, email)
}

// sanitizeEnvKey converts an email to a valid environment variable suffix.
// e.g., "user@example.com" -> "user_example_com"
func sanitizeEnvKey(email string) string {
	s := strings.ReplaceAll(email, "@", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ToUpper(s)
}
