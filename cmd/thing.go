package cmd

import (
	"fmt"
	"os"
	"piot-cli/api"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var (
	config_all        bool
	config_thing_type string
)

var thingCmd = &cobra.Command{
	Use:   "thing",
	Short: "Get list of things",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		things, err := client.GetThings(config_all)
		handleError(err)

		// use tabwriter.Debug flag (last arg) to see column borders
		w := tabwriter.NewWriter(os.Stdout, 0, 0, OUTPUT_PADDING, ' ', 0)
		fmt.Fprintf(w, "NAME\tALIAS\tTYPE\tENABLED\t%s\tINFLUXDB\tMYSQL\t\n", fmt.Sprintf(DefaultColor, "LAST SEEN"))
		for i := 0; i < len(things); i++ {

			tm := time.Unix(int64(things[i].LastSeen), 0)
			td := time.Now().Sub(tm).Truncate(time.Second)
			age := formatAge(td)

			if things[i].LastSeenInterval > 0 {
				time_diff := time.Now().Unix() - int64(things[i].LastSeen)
				if time_diff > int64(things[i].LastSeenInterval) {
					age = fmt.Sprintf(RedColor, age)
				} else {
					age = fmt.Sprintf(GreenColor, age)
				}
			} else {
				age = fmt.Sprintf(DefaultColor, age)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\t%v\t%v\t\n",
				things[i].Name,
				things[i].Alias,
				things[i].Type,
				things[i].Enabled,
				age,
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
	Short: "Create new thing",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		id, err := client.CreateThing(args[0], config_thing_type)
		handleError(err)

		log.Infof("%s", id)
	},
}

func init() {
	rootCmd.AddCommand(thingCmd)
	thingCmd.Flags().BoolVar(&config_all, "all", false, "Show all things across orgs")

	thingCmd.AddCommand(thingDeleteCmd)

	thingCmd.AddCommand(thingCreateCmd)
	thingCreateCmd.Flags().StringVar(&config_thing_type, "type", "device", "Thing type (device, sensor, switch)")
}
