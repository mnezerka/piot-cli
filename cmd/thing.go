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
	config_long       bool
)

var thingCmd = &cobra.Command{
	Use:   "thing",
	Short: "Get list of things",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		things, err := client.GetThings(config_all, nil)
		handleError(err)

		// use tabwriter.Debug flag (last arg) to see column borders
		w := tabwriter.NewWriter(os.Stdout, 0, 0, OUTPUT_PADDING, ' ', 0)
		if config_long {
			fmt.Fprintf(w, "ID\tNAME\tALIAS\tTYPE/CLASS\tENABLED\t%s\tVALUE\tINFLUXDB\tMYSQL\t\n", fmt.Sprintf(DefaultColor, "LAST SEEN"))
		} else {
			fmt.Fprintf(w, "NAME\tALIAS\tTYPE/CLASS\tENABLED\t%s\tVALUE\t\n", fmt.Sprintf(DefaultColor, "LAST SEEN"))
		}
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

			if config_long {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%s\t%v\t%v\t%v\t\n",
					things[i].Id,
					things[i].Name,
					things[i].Alias,
					things[i].Type+"/"+things[i].Sensor.Class,
					things[i].Enabled,
					age,
					things[i].Sensor.Value,
					things[i].StoreInfluxDb,
					things[i].StoreMysqlDb,
				)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\t%v\t\n",
					things[i].Name,
					things[i].Alias,
					things[i].Type+"/"+things[i].Sensor.Class,
					things[i].Enabled,
					age,
					things[i].Sensor.Value,
				)
			}
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
	thingCmd.Flags().BoolVarP(&config_long, "long", "l", false, "use long listing (show more columns)")

	thingCmd.AddCommand(thingDeleteCmd)

	thingCmd.AddCommand(thingCreateCmd)
	thingCreateCmd.Flags().StringVar(&config_thing_type, "type", "device", "Thing type (device, sensor, switch)")
}
