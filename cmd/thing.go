package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"text/tabwriter"
	"piot-cli/api"
	"time"
)

var thingCmd = &cobra.Command{
	Use:   "thing",
	Short: "Get list of things",
	Long: ``,
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

func init() {
	rootCmd.AddCommand(thingCmd)

}
