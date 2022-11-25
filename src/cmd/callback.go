package cmd

import (
	"github.com/idprm/go-yellowclinic/src/database"
	"github.com/idprm/go-yellowclinic/src/model"
	"github.com/spf13/cobra"
)

var callbackCmd = &cobra.Command{
	Use:   "callback",
	Short: "Callback CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var chat model.Chat
		database.Datasource.DB().Where("is_leave", false).Find(&chat)

	},
}

func init() {
	rootCmd.AddCommand(callbackCmd)
}
