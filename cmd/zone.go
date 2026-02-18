package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Manage Cloudflare DNS zones",
}

var zoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all zones in the account",
	RunE: func(cmd *cobra.Command, args []string) error {
		zones, err := client.ListZones()
		if err != nil {
			return err
		}

		if len(zones) == 0 {
			fmt.Println("No zones found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS")
		fmt.Fprintln(w, "──\t────\t──────")
		for _, z := range zones {
			fmt.Fprintf(w, "%s\t%s\t%s\n", z.ID, z.Name, z.Status)
		}
		w.Flush()

		return nil
	},
}

func init() {
	zoneCmd.AddCommand(zoneListCmd)
	rootCmd.AddCommand(zoneCmd)
}
