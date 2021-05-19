package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/op/go-logging"
)

const (
   OUTPUT_PADDING = 3
   LOGGER_MODULE = "piot"
   //LOGGER_FORMAT = "%{color}%{time:2006/01/02 15:04:05} [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"
   //LOGGER_FORMAT = "%{color}# [%{level:.6s}] %{shortfile} : %{color:reset}%{message}"
   LOGGER_FORMAT = "[%{level:.6s}] %{message}"
   LOGGER_FORMAT_COLORS = "%{color}[%{level:.6s}] %{color:reset}%{message}"
)

var (
	config_user string
	config_password string
	config_log_level string
)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "piot-cli",
	Short: "PIOT client",
	Long: ``,
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

	rootCmd.PersistentFlags().StringVar(&config_user, "user", "", "User")
	rootCmd.PersistentFlags().StringVar(&config_password, "password", "", "Password")
	rootCmd.PersistentFlags().StringVarP(&config_log_level, "log-level", "", "INFO", "Log level (CRITICIAL, ERROR, WARNING, NOTICE, INFO, DEBUG)")

	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetEnvPrefix("piot")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetDefault("piot.host", "iot.pavoucek.net/api")
	viper.AutomaticEnv() // read in environment variables that match

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
	logging.SetLevel(logLevel, LOGGER_MODULE);
}
