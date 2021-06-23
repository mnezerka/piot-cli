package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"piot-cli/api"
	"sort"
	"strings"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/jszwec/csvutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config_from  string
	config_to    string
	config_names string
)

const TIME_LAYOUT string = "2006-01-02"
const SENSOR_VALUE_EMPTY float64 = 9999

//const DATE_LAYOUT string = "2006-01-02T15:04:05"

type SensorValue struct {
	Date  time.Time `json:"date" csv:"date"`
	Value float64   `json:"value" csv:"value"`
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

	// response value can be nil in case there are no measurements in whole
	// grouping interval (1 hour)
	if response[1] == nil {
		result.Value = SENSOR_VALUE_EMPTY
	} else {
		result.Value, err = response[1].(json.Number).Float64()
		if err != nil {
			return nil, fmt.Errorf("Cannot parse sensor value from InfluxDB response (%v): %v", response, err)
		}
	}

	return &result, nil
}

func SensorData2CsvRows(sensor_data map[string][]SensorValue) (string, error) {

	// prepare commplete list of date time stamps (first column)
	// note: the map is not sorted !, it must be sorted e.g. before writing
	//       to output stream
	rows := map[time.Time][]float64{}

	for _, sensor_values := range sensor_data {
		for _, sensor_value := range sensor_values {
			if _, ok := rows[sensor_value.Date]; !ok {
				rows[sensor_value.Date] = []float64{}
			}
		}
	}

	header := []string{"date"}

	for sensor_name, sensor_values := range sensor_data {
		header = append(header, sensor_name)

		// go through global table rows
		for ts, _ := range rows {
			// add value for timestamp o empty value (0)
			exists := false
			for _, sensor_value := range sensor_values {
				if sensor_value.Date == ts {
					rows[ts] = append(rows[ts], sensor_value.Value)
					exists = true
					break
				}
			}
			if !exists {
				rows[ts] = append(rows[ts], SENSOR_VALUE_EMPTY)
			}
		}
	}

	// build csv
	var records [][]string

	// csv header
	records = append(records, header)

	time_stamps_sorted := make([]time.Time, 0, len(rows))
	for time_stamp := range rows {
		time_stamps_sorted = append(time_stamps_sorted, time_stamp)
	}
	sort.Slice(time_stamps_sorted, func(i, j int) bool {
		return time_stamps_sorted[i].Before(time_stamps_sorted[j])
	})

	// loop through rows in time sequence (sorted keys)
	for _, time_stamp := range time_stamps_sorted {
		values := rows[time_stamp]
		row_str := []string{time_stamp.String()}
		for _, value := range values {
			if value == SENSOR_VALUE_EMPTY {
				row_str = append(row_str, "nil")
			} else {
				row_str = append(row_str, fmt.Sprintf("%.2f", value))
			}
		}
		records = append(records, row_str)
	}

	// golang struct -> CSV string
	buf := bytes.NewBufferString("")
	w := csv.NewWriter(buf)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
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

		var names []string
		if config_names != "" {
			names = strings.Split(config_names, ",")
		}

		// TODO: check if to > from

		// get api client
		client := api.NewClient(log)
		err = client.Login()
		handleError(err)

		// get active org from user profile
		profile, err := client.GetUserProfile()
		handleError(err)
		org, err := profile.GetActiveOrg()
		handleError(err)

		log.Infof("Export params:")
		log.Infof("  from: %s", date_from)
		log.Infof("  to: %s", date_to)
		log.Infof("  names: %s", names)

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

		sensor_data := map[string][]SensorValue{}

		// fetch data for sensors (one by one)
		for _, thing := range things {

			// filter things if names flag was specified
			if len(names) > 0 {
				skip := true
				for _, name := range names {
					if thing.Name == name {
						skip = false
						break
					}
				}
				if skip {
					log.Infof("Skipping sensor '%s'", thing.Name)
					continue
				}
			}

			log.Infof("Fetching data for sensor '%s.%s'", org.Name, thing.Name)
			query := fmt.Sprintf(
				"SELECT MEAN(\"value\") FROM \"sensor\" WHERE time >= '%s' AND time <= '%s' AND \"id\" = '%s' GROUP BY time(1h)",
				date_from.Format(time.RFC3339),
				date_to.Format(time.RFC3339),
				thing.Id)

			log.Debugf("query: %s", query)

			q := influx.NewQuery(query, org.InfluxDb, "")

			response, err := ic.Query(q)
			handleError(err)

			if response.Error() != nil {
				handleError(response.Error())
			}

			log.Debugf("response: %s", response)

			if len(response.Results[0].Series) == 0 {
				log.Infof("No influxdb data for sensor  '%s.%s'", org.Name, thing.Name)
				sensor_data[thing.Name] = []SensorValue{}
				continue
			}

			sensor_data[thing.Name] = []SensorValue{}

			// we are interested in results from first statement and first entry from series
			for _, value := range response.Results[0].Series[0].Values {
				//fmt.Printf("%v\n", value)

				sensor_value, err := CreateFromInfluxResponse(value)
				handleError(err)

				sensor_data[thing.Name] = append(sensor_data[thing.Name], *sensor_value)
			}
		}

		switch config_format {
		case "csv":
			//result_csv, err := csvutil.Marshal(sensor_data)
			result_csv, err := SensorData2CsvRows(sensor_data)
			handleError(err)
			fmt.Println(string(result_csv))
		case "json", "":
			result_json, err := json.MarshalIndent(sensor_data, "", "  ")
			handleError(err)
			fmt.Println(string(result_json))
		default:
			handleError(fmt.Errorf("Unkonwn output format: %s", config_format))
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
	exportSensorsCmd.Flags().StringVar(&config_names, "names", "", "limit export to particular sensor names (comma seperated list)")
}
