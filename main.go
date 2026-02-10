package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"

	"github.com/andatoshiki/termfolio/auth"
	"github.com/andatoshiki/termfolio/config"
	"github.com/andatoshiki/termfolio/counter"
	"github.com/andatoshiki/termfolio/ui"
	"github.com/andatoshiki/termfolio/version"
)

func main() {
	// Parse CLI flags
	configPath := flag.String("c", "config.yaml", "Path to configuration file")
	showVersion := flag.Bool("v", false, "Show version information")
	flag.Parse()

	// Show version and exit if -v flag is set
	if *showVersion {
		fmt.Println(version.VersionInfo())
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

	var counterStore *counter.Store
	if cfg.Counter.Enabled {
		store, err := counter.Open(cfg.Counter.DBPath)
		if err != nil {
			log.Fatalf("Failed to open counter db: %v", err)
		}
		counterStore = store
	}

	// Ensure host key exists (will prompt user to generate if needed)
	if err := EnsureHostKey(cfg.SSH.HostKeyPath); err != nil {
		log.Fatalf("Failed to ensure host key: %v", err)
	}

	// Configure SSH public key authentication based on the auth mode
	publicKeyAuth := auth.PublicKeyHandler(cfg.SSH.AuthMode, cfg.SSH.AuthorizedKeys)

	teaHandler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		visitorCount := 0
		trackingEnabled := counterStore != nil
		remoteIP := ""
		if counterStore != nil {
			if addr := s.RemoteAddr(); addr != nil {
				host, _, err := net.SplitHostPort(addr.String())
				if err == nil {
					remoteIP = host
				} else {
					remoteIP = addr.String()
				}
			}

			if remoteIP != "" {
				optedOut, err := counterStore.IsOptedOut(remoteIP)
				if err != nil {
					log.Printf("Failed to read privacy status: %v", err)
				} else {
					trackingEnabled = !optedOut
				}
			}

			var err error
			if trackingEnabled {
				visitorCount, err = counterStore.RecordVisit(remoteIP)
			} else {
				visitorCount, err = counterStore.Count()
			}
			if err != nil {
				log.Printf("Failed to update counter: %v", err)
			}
		}
		return ui.NewModelWithCounter(counterStore, visitorCount, remoteIP, trackingEnabled), []tea.ProgramOption{tea.WithAltScreen()}
	}

	s, err := wish.NewServer(
		wish.WithAddress(cfg.SSH.ListenAddr()),
		wish.WithHostKeyPath(cfg.SSH.HostKeyPath),
		wish.WithPublicKeyAuth(publicKeyAuth),
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
