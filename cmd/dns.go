package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage Cloudflare DNS records",
}

var dnsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all DNS records for a zone",
	Example: `  gwaihir dns list --zone example.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, _ := cmd.Flags().GetString("zone")
		if zone == "" {
			return fmt.Errorf("--zone is required")
		}

		zoneID, err := client.GetZoneID(zone)
		if err != nil {
			return err
		}

		records, err := client.ListRecords(zoneID)
		if err != nil {
			return err
		}

		if len(records) == 0 {
			fmt.Printf("No records found for %s\n", zone)
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tNAME\tCONTENT\tTTL\tPROXIED")
		fmt.Fprintln(w, "──\t────\t────\t───────\t───\t───────")
		for _, r := range records {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%t\n", r.ID, r.Type, r.Name, r.Content, r.TTL, r.Proxied)
		}
		w.Flush()

		return nil
	},
}

func init() {
	dnsListCmd.Flags().String("zone", "", "Domain name (e.g. example.com)")

	dnsCmd.AddCommand(dnsListCmd)
	rootCmd.AddCommand(dnsCmd)
}
