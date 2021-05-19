package cmd

import (
	"fmt"
	"os"
	"piot-cli/api"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var thingCmd = &cobra.Command{
	Use:   "thing",
	Short: "Get list of things",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		things, err := client.GetThings(nil)
		handleError(err)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, OUTPUT_PADDING, ' ', 0)
		fmt.Fprintf(w, "ID\tNAME\tALIAS\tTYPE\tENABLED\tLAST SEEN\tINFLUXDB\tMYSQL\n")
		for i := 0; i < len(things); i++ {

			tm := time.Unix(int64(things[i].LastSeen), 0)

			td := time.Now().Sub(tm).Truncate(time.Second)

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%v\t%v\t%v\n",
				things[i].Id,
				things[i].Name,
				things[i].Alias,
				things[i].Type,
				things[i].Enabled,
				td,
				things[i].StoreInfluxDb,
				things[i].StoreMysqlDb,
			)
		}
		w.Flush()
	},
}

var thingDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete thing",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		err = client.DeleteThing()
		handleError(err)
	},
}

var thingCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create thing",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		err = client.CreateThing()
		handleError(err)
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(thingDeleteCmd)

	profileCmd.AddCommand(thingCreateCmd)
}
