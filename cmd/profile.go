package cmd

import (
	"encoding/json"
	"fmt"
	"piot-cli/api"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Get user profile",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient(log)

		err := client.Login()
		handleError(err)

		profile, err := client.GetUserProfile()
		handleError(err)

		profileJson, err := json.MarshalIndent(profile, "", "  ")
		handleError(err)

		fmt.Printf("%s\n", string(profileJson))
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
