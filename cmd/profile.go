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

//{"query":"mutation {updateUserProfile(profile: {org_id: \"5e1437163afe8695f1351311\"}) {is_admin, email, org_id, orgs {id, name}}}"}

func init() {
	rootCmd.AddCommand(profileCmd)
}
