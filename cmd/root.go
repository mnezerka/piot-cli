package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	OUTPUT_PADDING = 3
	LOGGER_MODULE  = "piot"
	//LOGGER_FORMAT = "%{color}%{time:2006/01/02 15:04:05} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"
	//LOGGER_FORMAT = "%{color}# [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"
	LOGGER_FORMAT        = "[%{level:.6s}] %{message}"
	LOGGER_FORMAT_COLORS = "%{color}[%{level:.6s}] %{color:reset}%{message}"
)

var (
	config_cfg_file          string
	config_piot_url          string
	config_piot_user         string
	config_piot_password     string
	config_log_level         string
	config_format            string
	config_influxdb_url      string
	config_influxdb_user     string
	config_influxdb_password string

//	config_org       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "piot-cli",
	Short:   "PIOT client",
	Long:    ``,
	Version: appVersion,
}

// global instance of logger
var log = logging.MustGetLogger(LOGGER_MODULE)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&config_cfg_file, "config", "", "config file (default is $HOME/.piot)")
	rootCmd.PersistentFlags().StringVar(&config_piot_url, "piot-url", "", "PIOT API url")
	rootCmd.PersistentFlags().StringVar(&config_piot_user, "piot-user", "", "User")
	rootCmd.PersistentFlags().StringVar(&config_piot_password, "piot-password", "", "Password")
	rootCmd.PersistentFlags().StringVarP(&config_log_level, "log-level", "", "INFO", "Log level (CRITICIAL, ERROR, WARNING, NOTICE, INFO, DEBUG)")
	//	rootCmd.PersistentFlags().StringVar(&config_org, "org", "", "Organization")

	rootCmd.PersistentFlags().StringVar(&config_influxdb_url, "influxdb-url", "", "InfluxDB URL")
	rootCmd.PersistentFlags().StringVar(&config_influxdb_user, "influxdb-user", "", "InfluxDB User")
	rootCmd.PersistentFlags().StringVar(&config_influxdb_password, "influxdb-password", "", "InfluxDB Password")

	viper.BindPFlag("piot.url", rootCmd.PersistentFlags().Lookup("piot-url"))
	viper.BindPFlag("piot.user", rootCmd.PersistentFlags().Lookup("piot-user"))
	viper.BindPFlag("piot.password", rootCmd.PersistentFlags().Lookup("piot-password"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("influxdb.url", rootCmd.PersistentFlags().Lookup("influxdb-url"))
	viper.BindPFlag("influxdb.user", rootCmd.PersistentFlags().Lookup("influxdb-user"))
	viper.BindPFlag("influxdb.password", rootCmd.PersistentFlags().Lookup("influxdb-password"))
	//	viper.BindPFlag("org", rootCmd.PersistentFlags().Lookup("org"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if config_cfg_file != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config_cfg_file)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".piot" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".piot")
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("piot")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	var configFileUsed string
	if err := viper.ReadInConfig(); err == nil {
		configFileUsed = viper.ConfigFileUsed()
	}

	// configure logging
	var logLevelStr = viper.GetString("log.level")
	// try to convert string log level
	logLevel, err := logging.LogLevel(logLevelStr)
	if err != nil {
		fmt.Printf("Invalid logging level: \"%s\"\n", logLevelStr)
		os.Exit(1)
	}

	formatterStdErr := logging.NewBackendFormatter(
		// out, prefix flag
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(LOGGER_FORMAT_COLORS),
	)
	logging.SetBackend(formatterStdErr)
	logging.SetLevel(logLevel, LOGGER_MODULE)

	if len(configFileUsed) > 0 {
		log.Infof("Using config file: '%s'", configFileUsed)
	}
}
