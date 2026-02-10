package auth

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/ssh"
)

// PublicKeyHandler creates an SSH public key authentication handler based on the auth mode
func PublicKeyHandler(authMode string, authorizedKeysPath string) func(ctx ssh.Context, key ssh.PublicKey) bool {
	switch authMode {
	case "none":
		// No public key authentication - reject all keys
		// This is the most secure option for a public-facing portfolio
		return func(ctx ssh.Context, key ssh.PublicKey) bool {
			return false
		}

	case "authorized_keys":
		// Load authorized keys and validate against them
		authorizedKeys, err := LoadAuthorizedKeys(authorizedKeysPath)
		if err != nil {
			log.Printf("WARNING: Failed to load authorized keys from %s: %v", authorizedKeysPath, err)
			log.Printf("WARNING: Rejecting all public key authentication attempts")
			return func(ctx ssh.Context, key ssh.PublicKey) bool {
				return false
			}
		}

		log.Printf("Loaded %d authorized keys from %s", len(authorizedKeys), authorizedKeysPath)

		return func(ctx ssh.Context, key ssh.PublicKey) bool {
			return IsKeyAuthorized(key, authorizedKeys)
		}

	case "allow_all":
		// INSECURE: Accept all public keys without validation
		log.Printf("WARNING: SSH authentication is set to 'allow_all' mode")
		log.Printf("WARNING: This accepts ANY SSH public key without validation")
		log.Printf("WARNING: This is insecure and should only be used for testing")
		log.Printf("WARNING: Consider using 'none' mode for public portfolios")

		return func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}

	default:
		// Unknown mode - reject all for safety
		log.Printf("ERROR: Unknown auth mode %q, rejecting all authentication", authMode)
		return func(ctx ssh.Context, key ssh.PublicKey) bool {
			return false
		}
	}
}

// LoadAuthorizedKeys loads SSH public keys from an OpenSSH authorized_keys formatted file
func LoadAuthorizedKeys(path string) ([]ssh.PublicKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open authorized keys file: %w", err)
	}
	defer file.Close()

	var keys []ssh.PublicKey
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the public key
		pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
		if err != nil {
			log.Printf("WARNING: Failed to parse key at line %d in %s: %v", lineNum, path, err)
			continue
		}

		keys = append(keys, pubKey)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading authorized keys file: %w", err)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no valid keys found in %s", path)
	}

	return keys, nil
}

// IsKeyAuthorized checks if a given public key is in the list of authorized keys
func IsKeyAuthorized(key ssh.PublicKey, authorizedKeys []ssh.PublicKey) bool {
	keyBytes := key.Marshal()

	for _, authKey := range authorizedKeys {
		if bytes.Equal(keyBytes, authKey.Marshal()) {
			return true
		}
	}

	return false
}
