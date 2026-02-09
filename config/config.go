package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SSH     SSHConfig     `yaml:"ssh"`
	Counter CounterConfig `yaml:"counter"`
}

type SSHConfig struct {
	Port        int    `yaml:"port"`
	Address     string `yaml:"address"`
	HostKeyPath string `yaml:"hostKeyPath"`
}

type CounterConfig struct {
	Enabled bool   `yaml:"enabled"`
	DBPath  string `yaml:"dbPath"`
}

func (c *CounterConfig) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		var enabled bool
		if err := value.Decode(&enabled); err != nil {
			return err
		}
		c.Enabled = enabled
		return nil
	case yaml.MappingNode:
		type counterYAML struct {
			Enabled *bool   `yaml:"enabled"`
			DBPath  *string `yaml:"dbPath"`
		}
		var raw counterYAML
		if err := value.Decode(&raw); err != nil {
			return err
		}
		if raw.Enabled != nil {
			c.Enabled = *raw.Enabled
		}
		if raw.DBPath != nil {
			c.DBPath = *raw.DBPath
		}
		return nil
	default:
		return fmt.Errorf("invalid counter config")
	}
}

// Load loads config from the specified path.
// If the file doesn't exist and userProvided is false, returns defaults.
// If the file doesn't exist and userProvided is true, returns an error.
func Load(configPath string, userProvided bool) (*Config, error) {
	// Start with defaults
	cfg := &Config{
		SSH: SSHConfig{
			Port:        2222,
			Address:     "0.0.0.0",
			HostKeyPath: ".ssh/host_ed25519",
		},
		Counter: CounterConfig{
			Enabled: true,
			DBPath:  "data/visitors.db",
		},
	}

	// Try to read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if userProvided {
				// User explicitly requested this file via -c flag—it must exist
				return nil, fmt.Errorf("config file not found at %s", configPath)
			}
			// Using default path and file doesn't exist—return defaults
			applyEnvVarOverrides(cfg)
			resolveHostKeyPath(cfg, configPath)
			resolveCounterPath(cfg, configPath)
			return cfg, nil
		}
		// Other read errors (permissions, etc.) are always errors
		return nil, fmt.Errorf("failed to read config file at %s: %w", configPath, err)
	}

	// Parse YAML and overlay onto defaults
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file at %s: %w", configPath, err)
	}

	// Environment variables override everything
	applyEnvVarOverrides(cfg)

	// Resolve host key path if relative (relative to config file directory)
	resolveHostKeyPath(cfg, configPath)
	resolveCounterPath(cfg, configPath)

	return cfg, nil
}

// applyEnvVarOverrides applies any set environment variables to the config,
// overriding both YAML and default values
func applyEnvVarOverrides(cfg *Config) {
	if port := os.Getenv("SSH_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil && p > 0 {
			cfg.SSH.Port = p
		}
	}

	if addr := os.Getenv("SSH_ADDRESS"); addr != "" {
		cfg.SSH.Address = addr
	}

	if hostKeyPath := os.Getenv("SSH_HOST_KEY_PATH"); hostKeyPath != "" {
		cfg.SSH.HostKeyPath = hostKeyPath
	}
}

func resolveHostKeyPath(cfg *Config, configPath string) {
	if cfg == nil {
		return
	}
	if cfg.SSH.HostKeyPath == "" || filepath.IsAbs(cfg.SSH.HostKeyPath) {
		return
	}
	if configPath == "" {
		return
	}
	baseDir := filepath.Dir(configPath)
	cfg.SSH.HostKeyPath = filepath.Clean(filepath.Join(baseDir, cfg.SSH.HostKeyPath))
}

func resolveCounterPath(cfg *Config, configPath string) {
	if cfg == nil {
		return
	}
	if cfg.Counter.DBPath == "" || filepath.IsAbs(cfg.Counter.DBPath) {
		return
	}
	if configPath == "" {
		return
	}
	baseDir := filepath.Dir(configPath)
	cfg.Counter.DBPath = filepath.Clean(filepath.Join(baseDir, cfg.Counter.DBPath))
}

func (cfg *SSHConfig) ListenAddr() string {
	return fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
}
