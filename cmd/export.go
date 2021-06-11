package cmd

import (
	"encoding/json"
	"fmt"
	"piot-cli/api"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/jszwec/csvutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config_from string
	config_to   string
)

const TIME_LAYOUT string = "2006-01-02"

//const DATE_LAYOUT string = "2006-01-02T15:04:05"

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

		things, err := client.GetThings(config_all, nil)
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
	Args: func(cmd *cobra.Command, args []string) error {
		_, err := time.Parse(TIME_LAYOUT, config_from)
		if err != nil {
			return err
		}
		_, err = time.Parse(TIME_LAYOUT, config_to)
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// convert from and to time.Time
		date_from, err := time.Parse(TIME_LAYOUT, config_from)
		handleError(err)
		date_to, err := time.Parse(TIME_LAYOUT, config_to)
		handleError(err)

		// get api client
		client := api.NewClient(log)
		err = client.Login()
		handleError(err)

		// get active org from user profile
		profile, err := client.GetUserProfile()
		handleError(err)
		org, err := profile.GetActiveOrg()
		handleError(err)

		log.Infof("Influx url: %s", viper.GetString("influxdb.url"))
		log.Infof("Influx user: %s", viper.GetString("influxdb.user"))
		log.Infof("Influx password: %s", viper.GetString("influxdb.password"))

		ic, err := influx.NewHTTPClient(influx.HTTPConfig{
			Addr:     viper.GetString("influxdb.url"),
			Username: viper.GetString("influxdb.user"),
			Password: viper.GetString("influxdb.password"),
		})
		handleError(err)
		defer ic.Close()

		// get all org sensors
		things, err := client.GetThings(false, func(thing *api.Thing) bool { return thing.Type == "sensor" })
		handleError(err)

		for i := 0; i < len(things); i++ {

			log.Infof("Fetching data for sensor '%s.%s'", org.Name, things[i].Name)
			query := fmt.Sprintf(
				"SELECT MEAN(\"value\") FROM \"sensor\" WHERE time >= '%s' AND time <= '%s' AND \"id\" = '%s' GROUP BY time(1h)",
				date_from.Format(time.RFC3339),
				date_to.Format(time.RFC3339),
				things[i].Id)

			fmt.Println(query)

			q := influx.NewQuery(query, org.InfluxDb, "")

			response, err := ic.Query(q)
			handleError(err)

			if response.Error() != nil {
				handleError(response.Error())
			}

			fmt.Println(response.Results)

		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	//exportCmd.Flags().BoolVar(&config_all, "all", false, "Show all things across orgs")

	exportCmd.AddCommand(exportThingsCmd)
	exportThingsCmd.Flags().StringVar(&config_format, "format", "json", "output format (json, csv)")

	exportCmd.AddCommand(exportSensorsCmd)
	exportSensorsCmd.Flags().StringVar(&config_format, "format", "json", "output format (json, csv)")
	exportSensorsCmd.Flags().StringVar(&config_from, "from", "", "starting date in format "+TIME_LAYOUT)
	exportSensorsCmd.Flags().StringVar(&config_to, "to", "", "end date in format "+TIME_LAYOUT)
	exportSensorsCmd.MarkFlagRequired("from")
	exportSensorsCmd.MarkFlagRequired("to")
}
