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

type SensorValue struct {
	Date  time.Time
	Value float64
}

func CreateFromInfluxResponse(response []interface{}) (*SensorValue, error) {

	var err error

	if len(response) < 2 {
		return nil, fmt.Errorf("Cannot decode sensor value from InfluxDB response (%v)", response)
	}

	result := SensorValue{}

	result.Date, err = time.Parse(time.RFC3339, response[0].(string))
	if err != nil {
		return nil, fmt.Errorf("Cannot parse sensor time from InfluxDB response (%v): %v", response, err)
	}
	result.Value, err = response[1].(json.Number).Float64()
	if err != nil {
		return nil, fmt.Errorf("Cannot parse sensor value from InfluxDB response (%v): %v", response, err)
	}

	return &result, nil
}

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
	Run: func(cmd *cobra.Command, args []string) {

		var err error

		date_to := time.Now()
		date_from := date_to.Add((-1 * 24) * time.Hour) // last day

		if config_from != "" {
			// convert from and to time.Time
			date_from, err = time.Parse(TIME_LAYOUT, config_from)
			handleError(err)
		}

		if config_to != "" {
			// convert from and to time.Time
			date_to, err = time.Parse(TIME_LAYOUT, config_to)
			handleError(err)
		}

		// check if to > from

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

		log.Infof("From: %s", date_from)
		log.Infof("To: %s", date_to)

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

		result_json := map[string][]SensorValue{}

		// fetch data for sensors (one by one)
		for _, thing := range things {

			log.Infof("Fetching data for sensor '%s.%s'", org.Name, thing.Name)
			query := fmt.Sprintf(
				"SELECT MEAN(\"value\") FROM \"sensor\" WHERE time >= '%s' AND time <= '%s' AND \"id\" = '%s' GROUP BY time(1h)",
				date_from.Format(time.RFC3339),
				date_to.Format(time.RFC3339),
				thing.Id)

			q := influx.NewQuery(query, org.InfluxDb, "")

			response, err := ic.Query(q)
			handleError(err)

			if response.Error() != nil {
				handleError(response.Error())
			}

			result_json[thing.Name] = []SensorValue{}

			// we are interested in results from first statement and first entry from series
			for _, value := range response.Results[0].Series[0].Values {
				fmt.Printf("%v\n", value)

				sensor_value, err := CreateFromInfluxResponse(value)
				handleError(err)

				result_json[thing.Name] = append(result_json[thing.Name], *sensor_value)
			}

		}

		fmt.Printf("\n\n%v\n", result_json)
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
	//exportSensorsCmd.MarkFlagRequired("from")
	//exportSensorsCmd.MarkFlagRequired("to")
}
