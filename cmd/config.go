package cmd

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// list of keys that point to secret values (passwords, keys, etc.)
var blackList = []string{
	"password",
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		settings := viper.AllSettings()

		for key, val := range settings {
			if contains(blackList, key) {
				val = "*****"
			}
			v := reflect.ValueOf(val)
			if v.Kind() == reflect.Map {
				fmt.Printf("%s:\n", key)
				if m, ok := val.(map[string]interface{}); ok {
					for key2, val2 := range m {
						if contains(blackList, key2) {
							val2 = "*****"
						}
						fmt.Printf("  %s: %v\n", key2, val2)
					}
				}
			} else {
				fmt.Printf("%s: %v\n", key, val)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
