package cmd

import (
	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Applies DB migrations for hukz.",
	Run: func(cmd *cobra.Command, args []string) {
		// setup database
		config.SetupDB()

		// apply database migrations
		model.Migration()
	},
}
