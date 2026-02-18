package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dukerupert/gwaihir/pkg/cloudflare"
	"github.com/spf13/cobra"
)

const defaultTTL = 600

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

var dnsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DNS record",
	Example: `  gwaihir dns create --zone example.com --name test.example.com --type A --content 1.2.3.4
  gwaihir dns create --zone example.com --name test.example.com --type CNAME --content target.com --proxied`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, _ := cmd.Flags().GetString("zone")
		name, _ := cmd.Flags().GetString("name")
		recordType, _ := cmd.Flags().GetString("type")
		content, _ := cmd.Flags().GetString("content")
		ttl, _ := cmd.Flags().GetInt("ttl")
		proxied, _ := cmd.Flags().GetBool("proxied")

		if zone == "" || name == "" || recordType == "" || content == "" {
			return fmt.Errorf("--zone, --name, --type, and --content are required")
		}

		zoneID, err := client.GetZoneID(zone)
		if err != nil {
			return err
		}

		record := cloudflare.DNSRecord{
			Type:    recordType,
			Name:    name,
			Content: content,
			TTL:     ttl,
			Proxied: proxied,
		}

		result, err := client.CreateRecord(zoneID, record)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created %s record %s → %s (ID: %s)\n", result.Type, result.Name, result.Content, result.ID)
		return nil
	},
}

var dnsEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a DNS record by ID",
	Example: `  gwaihir dns edit --zone example.com --id abc123 --type A --name test.example.com --content 5.6.7.8
  gwaihir dns edit --zone example.com --id abc123 --type A --name test.example.com --content 5.6.7.8 --proxied`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, _ := cmd.Flags().GetString("zone")
		id, _ := cmd.Flags().GetString("id")
		name, _ := cmd.Flags().GetString("name")
		recordType, _ := cmd.Flags().GetString("type")
		content, _ := cmd.Flags().GetString("content")
		ttl, _ := cmd.Flags().GetInt("ttl")
		proxied, _ := cmd.Flags().GetBool("proxied")

		if zone == "" || id == "" || name == "" || recordType == "" || content == "" {
			return fmt.Errorf("--zone, --id, --name, --type, and --content are required")
		}

		zoneID, err := client.GetZoneID(zone)
		if err != nil {
			return err
		}

		record := cloudflare.DNSRecord{
			Type:    recordType,
			Name:    name,
			Content: content,
			TTL:     ttl,
			Proxied: proxied,
		}

		result, err := client.EditRecord(zoneID, id, record)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Updated %s record %s → %s (ID: %s)\n", result.Type, result.Name, result.Content, result.ID)
		return nil
	},
}

var dnsDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a DNS record by ID",
	Example: `  gwaihir dns delete --zone example.com --id abc123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zone, _ := cmd.Flags().GetString("zone")
		id, _ := cmd.Flags().GetString("id")

		if zone == "" || id == "" {
			return fmt.Errorf("--zone and --id are required")
		}

		zoneID, err := client.GetZoneID(zone)
		if err != nil {
			return err
		}

		if err := client.DeleteRecord(zoneID, id); err != nil {
			return err
		}

		fmt.Printf("✓ Deleted record %s from %s\n", id, zone)
		return nil
	},
}

func init() {
	// list flags
	dnsListCmd.Flags().String("zone", "", "Domain name (e.g. example.com)")

	// create flags
	dnsCreateCmd.Flags().String("zone", "", "Domain name (e.g. example.com)")
	dnsCreateCmd.Flags().String("name", "", "Full record name (e.g. test.example.com)")
	dnsCreateCmd.Flags().String("type", "", "Record type (A, AAAA, CNAME, TXT, MX, etc.)")
	dnsCreateCmd.Flags().String("content", "", "Record content (e.g. IP address or target)")
	dnsCreateCmd.Flags().Int("ttl", defaultTTL, "Time to live in seconds (1 = automatic)")
	dnsCreateCmd.Flags().Bool("proxied", false, "Enable Cloudflare proxy")

	// edit flags
	dnsEditCmd.Flags().String("zone", "", "Domain name (e.g. example.com)")
	dnsEditCmd.Flags().String("id", "", "Record ID to edit (get from dns list)")
	dnsEditCmd.Flags().String("name", "", "Full record name (e.g. test.example.com)")
	dnsEditCmd.Flags().String("type", "", "Record type (A, AAAA, CNAME, TXT, MX, etc.)")
	dnsEditCmd.Flags().String("content", "", "Record content (e.g. IP address or target)")
	dnsEditCmd.Flags().Int("ttl", defaultTTL, "Time to live in seconds (1 = automatic)")
	dnsEditCmd.Flags().Bool("proxied", false, "Enable Cloudflare proxy")

	// delete flags
	dnsDeleteCmd.Flags().String("zone", "", "Domain name (e.g. example.com)")
	dnsDeleteCmd.Flags().String("id", "", "Record ID to delete (get from dns list)")

	dnsCmd.AddCommand(dnsListCmd)
	dnsCmd.AddCommand(dnsCreateCmd)
	dnsCmd.AddCommand(dnsEditCmd)
	dnsCmd.AddCommand(dnsDeleteCmd)
	rootCmd.AddCommand(dnsCmd)
}
