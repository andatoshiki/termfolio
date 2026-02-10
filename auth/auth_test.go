package auth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/ssh"
)

// These are real ed25519 public keys for testing purposes (no private keys)
const (
	testKey1 = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG2kXStRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjR test1@example.com"
	testKey2 = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIF2kXStRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjR test2@example.com"
)

// mockPublicKey creates a mock SSH public key for testing
func mockPublicKey(data string) ssh.PublicKey {
	if data == "" {
		data = testKey1
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(data))
	if err != nil {
		// If parsing fails, we can't proceed with the test
		panic("failed to create mock key: " + err.Error())
	}
	return pubKey
}

func TestPublicKeyHandler_None(t *testing.T) {
	handler := PublicKeyHandler("none", "")

	// Should reject all keys
	key := mockPublicKey("")
	if handler(nil, key) {
		t.Error("Expected 'none' mode to reject keys, but it accepted")
	}
}

func TestPublicKeyHandler_AllowAll(t *testing.T) {
	handler := PublicKeyHandler("allow_all", "")

	// Should accept all keys
	key := mockPublicKey("")
	if !handler(nil, key) {
		t.Error("Expected 'allow_all' mode to accept keys, but it rejected")
	}
}

func TestPublicKeyHandler_UnknownMode(t *testing.T) {
	handler := PublicKeyHandler("invalid_mode", "")

	// Should reject all keys for unknown modes (fail-safe)
	key := mockPublicKey("")
	if handler(nil, key) {
		t.Error("Expected unknown mode to reject keys, but it accepted")
	}
}

func TestPublicKeyHandler_AuthorizedKeys_FileNotFound(t *testing.T) {
	handler := PublicKeyHandler("authorized_keys", "/nonexistent/path/authorized_keys")

	// Should reject all keys when file doesn't exist
	key := mockPublicKey("")
	if handler(nil, key) {
		t.Error("Expected to reject keys when authorized_keys file doesn't exist")
	}
}

func TestLoadAuthorizedKeys_ValidFile(t *testing.T) {
	// Create a temporary authorized_keys file
	tmpDir := t.TempDir()
	authKeysPath := filepath.Join(tmpDir, "authorized_keys")

	content := `# Comment line
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG2kXStRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjR test1@example.com
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIF2kXStRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjR test2@example.com
`

	if err := os.WriteFile(authKeysPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	keys, err := LoadAuthorizedKeys(authKeysPath)
	if err != nil {
		t.Fatalf("LoadAuthorizedKeys failed: %v", err)
	}

	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
}

func TestLoadAuthorizedKeys_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	authKeysPath := filepath.Join(tmpDir, "authorized_keys")

	// Create empty file
	if err := os.WriteFile(authKeysPath, []byte("# Just comments\n\n"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := LoadAuthorizedKeys(authKeysPath)
	if err == nil {
		t.Error("Expected error for empty authorized_keys file")
	}
}

func TestLoadAuthorizedKeys_InvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	authKeysPath := filepath.Join(tmpDir, "authorized_keys")

	content := `invalid-key-format
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG2kXStRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjRGkjR test@example.com
`

	if err := os.WriteFile(authKeysPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	keys, err := LoadAuthorizedKeys(authKeysPath)
	if err != nil {
		t.Fatalf("LoadAuthorizedKeys failed: %v", err)
	}

	// Should skip invalid line and only load valid key
	if len(keys) != 1 {
		t.Errorf("Expected 1 valid key, got %d", len(keys))
	}
}

func TestIsKeyAuthorized(t *testing.T) {
	key1 := mockPublicKey(testKey1)
	key2 := mockPublicKey(testKey2)

	authorizedKeys := []ssh.PublicKey{key1}

	// key1 should be authorized
	if !IsKeyAuthorized(key1, authorizedKeys) {
		t.Error("Expected key1 to be authorized")
	}

	// key2 should not be authorized
	if IsKeyAuthorized(key2, authorizedKeys) {
		t.Error("Expected key2 to not be authorized")
	}
}

func TestIsKeyAuthorized_EmptyList(t *testing.T) {
	key := mockPublicKey("")
	authorizedKeys := []ssh.PublicKey{}

	// Should not be authorized when list is empty
	if IsKeyAuthorized(key, authorizedKeys) {
		t.Error("Expected key to not be authorized with empty list")
	}
}

func TestPublicKeyHandler_AuthorizedKeys_Integration(t *testing.T) {
	// Create a temporary authorized_keys file
	tmpDir := t.TempDir()
	authKeysPath := filepath.Join(tmpDir, "authorized_keys")

	content := testKey1 + "\n"

	if err := os.WriteFile(authKeysPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	handler := PublicKeyHandler("authorized_keys", authKeysPath)

	// Create the exact same key that's in authorized_keys
	authorizedKey := mockPublicKey(testKey1)

	// Should accept the authorized key
	if !handler(nil, authorizedKey) {
		t.Error("Expected handler to accept authorized key")
	}

	// Create a different key
	unauthorizedKey := mockPublicKey(testKey2)

	// Should reject unauthorized key
	if handler(nil, unauthorizedKey) {
		t.Error("Expected handler to reject unauthorized key")
	}
}
