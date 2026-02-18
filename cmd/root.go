package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dukerupert/gwaihir/pkg/cloudflare"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var client *cloudflare.Client

var rootCmd = &cobra.Command{
	Use:           "gwaihir",
	Short:         "A CLI tool for managing Cloudflare DNS records",
	Long:          `Gwaihir manages Cloudflare DNS records via the Cloudflare API.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initClient()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initClient() error {
	// Try ~/.dotfiles/.env first, fall back to .env in current directory
	dotfilePath := filepath.Join(os.Getenv("HOME"), ".dotfiles", ".env")
	if _, err := os.Stat(dotfilePath); err == nil {
		_ = godotenv.Load(dotfilePath)
	} else {
		_ = godotenv.Load(".env")
	}

	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		return fmt.Errorf("CLOUDFLARE_API_TOKEN must be set in ~/.dotfiles/.env or .env")
	}

	client = cloudflare.NewClient(token)
	return nil
}
