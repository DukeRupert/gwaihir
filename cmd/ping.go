package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Verify API token and connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.VerifyToken(); err != nil {
			return err
		}
		fmt.Println("✓ Token is valid")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
