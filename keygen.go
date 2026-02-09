package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EnsureHostKey generates a host key if it doesn't exist
func EnsureHostKey(keyPath string) error {
	// Check if key exists
	if _, err := os.Stat(keyPath); err == nil {
		// Key exists
		return nil
	} else if !os.IsNotExist(err) {
		// Some other error occurred
		return err
	}

	// Key doesn't exist, prompt user
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("SSH host key not found at %s\n", keyPath)
	fmt.Print("Generate a new SSH host key? (y/n): ")

	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Aborted. Host key is required to run the SSH server.")
		os.Exit(1)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(keyPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Generate key using ssh-keygen
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate SSH key: %w", err)
	}

	fmt.Printf("SSH host key generated at %s\n", keyPath)
	return nil
}
