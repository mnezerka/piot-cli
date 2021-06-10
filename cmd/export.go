package cmd

import (
	"encoding/json"
	"fmt"
	"piot-cli/api"

	"github.com/jszwec/csvutil"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export piot data (sensors, things, etc.)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var exportThingsCmd = &cobra.Command{
	Use:   "things",
	Short: "Export things form current organization",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		things, err := client.GetThings(config_all)
		handleError(err)

		if config_format == "csv" {

			b, err := csvutil.Marshal(things)
			handleError(err)
			fmt.Println(string(b))

		} else {
			thingsJson, err := json.MarshalIndent(things, "", "  ")
			handleError(err)
			fmt.Printf("%s\n", string(thingsJson))
		}
	},
}

var exportSensorsCmd = &cobra.Command{
	Use:   "sensors",
	Short: "Export selected sensors",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	//exportCmd.Flags().BoolVar(&config_all, "all", false, "Show all things across orgs")

	exportCmd.AddCommand(exportThingsCmd)
	exportThingsCmd.Flags().StringVar(&config_format, "format", "json", "output format (json, csv)")

	exportCmd.AddCommand(exportSensorsCmd)
	exportSensorsCmd.Flags().StringVar(&config_format, "format", "json", "output format (json, csv)")
}
