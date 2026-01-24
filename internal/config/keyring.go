package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/zalando/go-keyring"
)

const serviceName = "sog"

// Protocol represents a service protocol for password lookup.
type Protocol string

const (
	ProtocolDefault Protocol = ""
	ProtocolIMAP    Protocol = "imap"
	ProtocolSMTP    Protocol = "smtp"
	ProtocolCalDAV  Protocol = "caldav"
	ProtocolCardDAV Protocol = "carddav"
	ProtocolWebDAV  Protocol = "webdav"
)

// SetPassword stores a password in the system keyring.
func SetPassword(email, password string) error {
	return keyring.Set(serviceName, email, password)
}

// SetPasswordForProtocol stores a protocol-specific password in the system keyring.
func SetPasswordForProtocol(email string, protocol Protocol, password string) error {
	if protocol == ProtocolDefault {
		return SetPassword(email, password)
	}
	key := fmt.Sprintf("%s:%s", email, protocol)
	return keyring.Set(serviceName, key, password)
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

// GetPasswordForProtocol retrieves a protocol-specific password.
// Falls back to: protocol-specific key → default key → environment variable.
func GetPasswordForProtocol(email string, protocol Protocol) (string, error) {
	// Try protocol-specific keyring key first
	if protocol != ProtocolDefault {
		key := fmt.Sprintf("%s:%s", email, protocol)
		password, err := keyring.Get(serviceName, key)
		if err == nil {
			return password, nil
		}

		// Try protocol-specific environment variable
		envKey := fmt.Sprintf("SOG_PASSWORD_%s_%s", sanitizeEnvKey(email), strings.ToUpper(string(protocol)))
		if envPass := os.Getenv(envKey); envPass != "" {
			return envPass, nil
		}
	}

	// Fall back to default password
	return GetPassword(email)
}

// DeletePassword removes a password from the system keyring.
func DeletePassword(email string) error {
	return keyring.Delete(serviceName, email)
}

// DeletePasswordForProtocol removes a protocol-specific password from the keyring.
func DeletePasswordForProtocol(email string, protocol Protocol) error {
	if protocol == ProtocolDefault {
		return DeletePassword(email)
	}
	key := fmt.Sprintf("%s:%s", email, protocol)
	return keyring.Delete(serviceName, key)
}

// sanitizeEnvKey converts an email to a valid environment variable suffix.
// e.g., "user@example.com" -> "user_example_com"
func sanitizeEnvKey(email string) string {
	s := strings.ReplaceAll(email, "@", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ToUpper(s)
}
