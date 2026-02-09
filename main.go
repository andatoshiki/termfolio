package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"

	"github.com/andatoshiki/termfolio/config"
	"github.com/andatoshiki/termfolio/ui"
)

func main() {
	// Parse CLI flags
	configPath := flag.String("c", "config.yaml", "Path to configuration file")
	showVersion := flag.Bool("v", false, "Show version information")
	flag.Parse()

	// Show version and exit if -v flag is set
	if *showVersion {
		fmt.Println(VersionInfo())
		os.Exit(0)
	}

	// Check if -c flag was explicitly provided
	userProvidedPath := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "c" {
			userProvidedPath = true
		}
	})

	// Set color profile
	lipgloss.SetColorProfile(termenv.ANSI256)

	// Load configuration
	cfg, err := config.Load(*configPath, userProvidedPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Ensure host key exists (will prompt user to generate if needed)
	if err := EnsureHostKey(cfg.SSH.HostKeyPath); err != nil {
		log.Fatalf("Failed to ensure host key: %v", err)
	}

	publicKeyAuth := func(ctx ssh.Context, key ssh.PublicKey) bool {
		return true
	}
	// This lets people see the portfolio even without a public key
	passwordAuth := func(ctx ssh.Context, password string) bool {
		return true
	}

	teaHandler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		return ui.NewModel(), []tea.ProgramOption{tea.WithAltScreen()}
	}

	s, err := wish.NewServer(
		wish.WithAddress(cfg.SSH.ListenAddr()),
		wish.WithHostKeyPath(cfg.SSH.HostKeyPath),
		wish.WithPublicKeyAuth(publicKeyAuth),
		wish.WithPasswordAuth(passwordAuth),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}
