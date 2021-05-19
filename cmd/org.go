package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"text/tabwriter"
	"piot-cli/api"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Get list of organizations",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		orgs, err := client.GetOrgs()
		handleError(err)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, OUTPUT_PADDING, ' ', 0)
		fmt.Fprintf(w, "ID\tNAME\n")
		for i := 0; i < len(orgs); i++ {

			fmt.Fprintf(w, "%s\t%s\n",
				orgs[i].Id,
				orgs[i].Name,
			)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
}
