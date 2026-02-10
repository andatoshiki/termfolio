package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateAuthMode_Valid(t *testing.T) {
	testCases := []struct {
		name     string
		authMode string
	}{
		{"none mode", "none"},
		{"allow_all mode", "allow_all"},
		{"authorized_keys mode with path", "authorized_keys"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{
				SSH: SSHConfig{
					AuthMode:       tc.authMode,
					AuthorizedKeys: "/path/to/keys",
				},
			}
			if err := validateAuthMode(cfg); err != nil {
				t.Errorf("Expected no error for %s, got: %v", tc.authMode, err)
			}
		})
	}
}

func TestValidateAuthMode_Empty(t *testing.T) {
	cfg := &Config{
		SSH: SSHConfig{
			AuthMode: "",
		},
	}
	if err := validateAuthMode(cfg); err != nil {
		t.Errorf("Expected no error for empty mode, got: %v", err)
	}
	// Should default to "none"
	if cfg.SSH.AuthMode != "none" {
		t.Errorf("Expected empty mode to default to 'none', got: %s", cfg.SSH.AuthMode)
	}
}

func TestValidateAuthMode_Invalid(t *testing.T) {
	cfg := &Config{
		SSH: SSHConfig{
			AuthMode: "invalid_mode",
		},
	}
	err := validateAuthMode(cfg)
	if err == nil {
		t.Error("Expected error for invalid auth mode")
	}
}

func TestValidateAuthMode_AuthorizedKeysWithoutPath(t *testing.T) {
	cfg := &Config{
		SSH: SSHConfig{
			AuthMode:       "authorized_keys",
			AuthorizedKeys: "",
		},
	}
	err := validateAuthMode(cfg)
	if err == nil {
		t.Error("Expected error when authorized_keys mode has no path")
	}
}

func TestLoad_DefaultConfig(t *testing.T) {
	// Load without existing file (should use defaults)
	cfg, err := Load("/nonexistent/config.yaml", false)
	if err != nil {
		t.Fatalf("Expected no error with default config, got: %v", err)
	}

	if cfg.SSH.Port != 2222 {
		t.Errorf("Expected default port 2222, got %d", cfg.SSH.Port)
	}
	if cfg.SSH.Address != "0.0.0.0" {
		t.Errorf("Expected default address 0.0.0.0, got %s", cfg.SSH.Address)
	}
	if cfg.SSH.AuthMode != "none" {
		t.Errorf("Expected default auth mode 'none', got %s", cfg.SSH.AuthMode)
	}
}

func TestLoad_WithAuthConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	content := `ssh:
  port: 3333
  address: "127.0.0.1"
  hostKeyPath: "test.key"
  authMode: "authorized_keys"
  authorizedKeys: "auth.keys"

counter:
  enabled: false
`

	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := Load(configPath, true)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.SSH.Port != 3333 {
		t.Errorf("Expected port 3333, got %d", cfg.SSH.Port)
	}
	if cfg.SSH.AuthMode != "authorized_keys" {
		t.Errorf("Expected auth mode 'authorized_keys', got %s", cfg.SSH.AuthMode)
	}
	if cfg.SSH.AuthorizedKeys == "" {
		t.Error("Expected authorizedKeys path to be set")
	}
}

func TestLoad_EnvVarOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("SSH_AUTH_MODE", "allow_all")
	os.Setenv("SSH_AUTHORIZED_KEYS", "/custom/path/keys")
	defer func() {
		os.Unsetenv("SSH_AUTH_MODE")
		os.Unsetenv("SSH_AUTHORIZED_KEYS")
	}()

	cfg, err := Load("/nonexistent/config.yaml", false)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if cfg.SSH.AuthMode != "allow_all" {
		t.Errorf("Expected auth mode from env 'allow_all', got %s", cfg.SSH.AuthMode)
	}
	if cfg.SSH.AuthorizedKeys != "/custom/path/keys" {
		t.Errorf("Expected authorized keys from env, got %s", cfg.SSH.AuthorizedKeys)
	}
}

func TestResolveAuthorizedKeysPath_Relative(t *testing.T) {
	cfg := &Config{
		SSH: SSHConfig{
			AuthorizedKeys: ".ssh/authorized_keys",
		},
	}

	configPath := "/home/user/config.yaml"
	resolveAuthorizedKeysPath(cfg, configPath)

	expected := filepath.Clean("/home/user/.ssh/authorized_keys")
	if cfg.SSH.AuthorizedKeys != expected {
		t.Errorf("Expected %s, got %s", expected, cfg.SSH.AuthorizedKeys)
	}
}

func TestResolveAuthorizedKeysPath_Absolute(t *testing.T) {
	absolutePath := "/etc/ssh/authorized_keys"
	cfg := &Config{
		SSH: SSHConfig{
			AuthorizedKeys: absolutePath,
		},
	}

	resolveAuthorizedKeysPath(cfg, "/home/user/config.yaml")

	// Should remain unchanged
	if cfg.SSH.AuthorizedKeys != absolutePath {
		t.Errorf("Expected absolute path unchanged, got %s", cfg.SSH.AuthorizedKeys)
	}
}
