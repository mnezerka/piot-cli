package cmd

import (
	"fmt"
	"os"
	"piot-cli/api"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Get list of organizations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		profile, err := client.GetUserProfile()
		handleError(err)

		orgs, err := client.GetOrgs(nil)
		handleError(err)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, OUTPUT_PADDING, ' ', 0)
		fmt.Fprintf(w, "NAME\tMEMBER\tCURRENT\t\n")
		for i := 0; i < len(orgs); i++ {

			isCurrent := ""
			if orgs[i].Id == profile.OrgId {
				isCurrent = "X"
			}

			isMember := ""
			for j := 0; j < len(profile.Orgs); j++ {
				if profile.Orgs[j].Id == orgs[i].Id {
					isMember = "X"
					break
				}
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t\n",
				orgs[i].Name,
				isMember,
				isCurrent,
			)
		}
		w.Flush()
	},
}

var orgSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set current org",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		err = client.SetCurrentOrg(args[0])
		handleError(err)
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)

	orgCmd.AddCommand(orgSetCmd)
	//orgSetCmd.Flags().StringVar(&config_org, "org", "", "Organization")
	//orgSetCmd.MarkFlagRequired("org")
}
