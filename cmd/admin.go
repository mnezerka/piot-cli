package cmd

import (
	"fmt"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Administration of PIOT environment",
	Long:  ``,
}

var adminInfluxDb = &cobra.Command{
	Use:   "influxdb",
	Short: "Administration of InfluxDb",
	Run: func(cmd *cobra.Command, args []string) {

		var err error

		ic, err := influx.NewHTTPClient(influx.HTTPConfig{
			Addr:     viper.GetString("influxdb.url"),
			Username: viper.GetString("influxdb.user"),
			Password: viper.GetString("influxdb.password"),
		})
		handleError(err)
		defer ic.Close()

		query := "SHOW DATABASES"

		log.Debugf("query: %s", query)

		q := influx.NewQuery(query, "", "")

		response, err := ic.Query(q)
		handleError(err)

		if response.Error() != nil {
			handleError(response.Error())
		}

		log.Debugf("response: %s", response)

		for _, db_name := range response.Results[0].Series[0].Values {
			fmt.Printf("%v\n", db_name[0])
		}
	},
}

var adminInfluxDbCreate = &cobra.Command{
	Use:   "create",
	Short: "Create database",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var err error

		ic, err := influx.NewHTTPClient(influx.HTTPConfig{
			Addr:     viper.GetString("influxdb.url"),
			Username: viper.GetString("influxdb.user"),
			Password: viper.GetString("influxdb.password"),
		})
		handleError(err)
		defer ic.Close()

		query := fmt.Sprintf("CREATE DATABASE \"%s\"", args[0])

		log.Debugf("query: %s", query)

		q := influx.NewQuery(query, "", "")

		response, err := ic.Query(q)
		handleError(err)

		if response.Error() != nil {
			handleError(response.Error())
		}

		log.Debugf("response: %s", response)
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)

	adminCmd.AddCommand(adminInfluxDb)

	adminInfluxDb.AddCommand(adminInfluxDbCreate)
}
